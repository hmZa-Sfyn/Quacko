// Import Go package
import "crypto/sha256"
import "encoding/hex"

// Use Go functions
fn hashString(input) {
    // Create hasher
    hasher = sha256.New()
    
    // Write data
    hasher.Write(bytes(input))
    
    // Get hash
    hash = hasher.Sum(nil)
    
    // Convert to hex string
    return hex.EncodeToString(hash)
}

password = "secret123"
hashedPassword = hashString(password)
println(hashedPassword)

// Register custom Go function for use in Quacko
go.Register("formatDate", fn(date, format) {
    // Go implementation
    return date.Format(format)
})

// Use the registered function
today = datetime.now()
formatted = formatDate(today, "2006-01-02")
println(formatted)