package utils

import (
        "strings"

        "github.com/gin-gonic/gin"
)

// GetCookieDomain determines the appropriate cookie domain based on the request host
// and the list of allowed domains.
func GetCookieDomain(c *gin.Context, allowedDomains []string) string {
        // If no domains are configured, return empty string
        // This tells the browser to use the current domain (no domain restriction)
        if len(allowedDomains) == 0 {
                return ""
        }

        host := c.Request.Host
        // Remove port if present
        if idx := strings.Index(host, ":"); idx != -1 {
                host = host[:idx]
        }

        for _, domain := range allowedDomains {
                if domain == host {
                        return domain
                }
        }

        // Return the first allowed domain as fallback
        return allowedDomains[0]
}
