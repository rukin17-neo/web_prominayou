package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
	// Critical auth endpoints - strict limits
	authLoginLimiter  *limiter.Limiter
	authForgotLimiter *limiter.Limiter
	authResetLimiter  *limiter.Limiter

	// High priority - user management
	userManagementLimiter *limiter.Limiter

	// Medium priority - normal CRUD
	crudLimiter *limiter.Limiter

	// Store for cleanup
	memoryStore limiter.Store
)

// InitRateLimiters initializes all rate limiters
func InitRateLimiters() error {
	// Create memory store with cleanup
	memoryStore = memory.NewStore()

	// Critical auth limits (IP-based)
	authLoginLimiter = limiter.New(memoryStore, limiter.Rate{
		Period: 15 * time.Minute,
		Limit:  5,
	})

	authForgotLimiter = limiter.New(memoryStore, limiter.Rate{
		Period: 15 * time.Minute,
		Limit:  3,
	})

	authResetLimiter = limiter.New(memoryStore, limiter.Rate{
		Period: 15 * time.Minute,
		Limit:  5,
	})

	// High priority (session/IP-based)
	userManagementLimiter = limiter.New(memoryStore, limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  20,
	})

	// Medium priority CRUD
	crudLimiter = limiter.New(memoryStore, limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  30,
	})

	log.Println("Rate limiters initialized successfully")
	return nil
}

// RateLimitLogin creates rate limit middleware for login
func RateLimitLogin() func(http.Handler) http.Handler {
	return createRateLimiter(authLoginLimiter, "login")
}

// RateLimitForgotPassword creates rate limit middleware for forgot password
func RateLimitForgotPassword() func(http.Handler) http.Handler {
	return createRateLimiter(authForgotLimiter, "forgot-password")
}

// RateLimitResetPassword creates rate limit middleware for reset password
func RateLimitResetPassword() func(http.Handler) http.Handler {
	return createRateLimiter(authResetLimiter, "reset-password")
}

// RateLimitUserManagement creates rate limit middleware for user management
func RateLimitUserManagement() func(http.Handler) http.Handler {
	return createRateLimiter(userManagementLimiter, "user-management")
}

// RateLimitCRUD creates rate limit middleware for general CRUD operations
func RateLimitCRUD() func(http.Handler) http.Handler {
	return createRateLimiter(crudLimiter, "crud")
}

// createRateLimiter is a helper function to create rate limit middleware
func createRateLimiter(lim *limiter.Limiter, name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only apply rate limiting to POST requests
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}

			ip := getIP(r)
			context, err := lim.Get(r.Context(), ip)
			if err != nil {
				log.Printf("Rate limiter error for %s: %v", name, err)
				// On error, allow the request to proceed (fail open)
				next.ServeHTTP(w, r)
				return
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

			if context.Reached {
				resetTime := time.Unix(context.Reset, 0)
				retryAfter := int(time.Until(resetTime).Seconds())

				log.Printf("RATE_LIMIT_EXCEEDED: endpoint=%s, ip=%s, limit=%d, window=%s",
					r.URL.Path, ip, context.Limit, name)

				w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
				http.Error(w, fmt.Sprintf("Too many requests. Please try again in %d seconds.", retryAfter),
					http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
