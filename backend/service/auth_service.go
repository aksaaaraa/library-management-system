package service

func (s *AuthService) RefreshToken(refreshToken, ipAddress, userAgent string) (*models.TokenPair, *models.User, error) {
	// Verify that the token is a refresh token
	if token.Claims.(jwt.MapClaims)["type"] != "refresh" {
	return nil, nil, models.ErrInvalidToken
}

// Get token hash
tokenHash := token.Claims.(jwt.MapClaims)["jti"].(string)

// Check if token exists in the database
storedToken, err := s.tokenRepo.GetByHash(tokenHash)
if err != nil || storedToken == nil {
	s.authLogRepo.Create(&models.AuthLog{
		Action:    models.AuthActionRefreshToken,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    models.StatusFailure,
		Details:   "Token not found in database",
	})
	return nil, nil, models.ErrInvalidToken
}

// Check if token is expired
if storedToken.ExpiresAt.Before(time.Now()) {
	s.authLogRepo.Create(&models.AuthLog{
		UserID:    storedToken.UserID,
		Action:    models.AuthActionRefreshToken,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    models.StatusFailure,
		Details:   "Token expired",
	})
	return nil, nil, models.ErrTokenExpired
}

// Check if token type is refresh
if storedToken.TokenType != models.TokenTypeRefresh {
	s.authLogRepo.Create(&models.AuthLog{
		UserID:    storedToken.UserID,
		Action:    models.AuthActionRefreshToken,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    models.StatusFailure,
		Details:   "Not a refresh token",
	})
	return nil, nil, models.ErrInvalidToken
}

// Get user from database
user, err := s.userRepo.GetByID(storedToken.UserID)
if err != nil {
	return nil, nil, fmt.Errorf("failed to get user: %w", err)
}

// Check account status
if user.AccountStatus != models.UserStatusActive {
	s.authLogRepo.Create(&models.AuthLog{
		UserID:    user.ID,
		Action:    models.AuthActionRefreshToken,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    models.StatusFailure,
		Details:   fmt.Sprintf("Account status: %s", user.AccountStatus),
	})
	return nil, nil, models.ErrAccountNotActive
}

// Revoke the old refresh token
s.tokenRepo.Delete(storedToken.ID)

// Generate new token pair
tokenPair, err := s.generateTokenPair(user)
if err != nil {
	return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
}

// Store new tokens
accessTokenExpires := time.Now().Add(time.Duration(s.config.AccessTokenLifetimeMinutes) * time.Minute)
refreshTokenExpires := time.Now().Add(time.Duration(s.config.RefreshTokenLifetimeDays) * 24 * time.Hour)

err = s.tokenRepo.Create(&models.AccessToken{
	UserID:    user.ID,
	TokenHash: tokenPair.AccessTokenHash,
	TokenType: models.TokenTypeAccess,
	ExpiresAt: accessTokenExpires,
	IPAddress: ipAddress,
	UserAgent: userAgent,
})
if err != nil {
	return nil, nil, fmt.Errorf("failed to store access token: %w", err)
}

err = s.tokenRepo.Create(&models.AccessToken{
	UserID:    user.ID,
	TokenHash: tokenPair.RefreshTokenHash,
	TokenType: models.TokenTypeRefresh,
	ExpiresAt: refreshTokenExpires,
	IPAddress: ipAddress,
	UserAgent: userAgent,
})
if err != nil {
	return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
}

// Log successful token refresh
s.authLogRepo.Create(&models.AuthLog{
	UserID:    user.ID,
	Action:    models.AuthActionRefreshToken,
	IPAddress: ipAddress,
	UserAgent: userAgent,
	Status:    models.StatusSuccess,
})

return tokenPair, user, nil
}

return nil, nil, fmt.Errorf("invalid token"){
}

// Logout invalidates a user's tokens
func (s *AuthService) Logout(accessToken, ipAddress, userAgent string) error {
// Parse and validate access token
token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
}
return []byte(s.config.JWTSecret), nil
})

if err != nil || !token.Valid {
return models.ErrInvalidToken
}

claims, ok := token.Claims.(jwt.MapClaims)
if !ok {
return models.ErrInvalidToken
}

userID, ok := claims["sub"].(string)
if !ok {
return models.ErrInvalidToken
}

tokenHash, ok := claims["jti"].(string)
if !ok {
return models.ErrInvalidToken
}

// Find and delete the access token
storedToken, err := s.tokenRepo.GetByHash(tokenHash)
if err == nil && storedToken != nil {
s.tokenRepo.Delete(storedToken.ID)
}

// Log the logout
s.authLogRepo.Create(&models.AuthLog{
UserID:    userID,
Action:    models.AuthActionLogout,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
})

return nil
}

// LogoutAllSessions invalidates all tokens for a specific user
func (s *AuthService) LogoutAllSessions(userID, ipAddress, userAgent string) error {
// Delete all tokens for this user
err := s.tokenRepo.DeleteAllForUser(userID)
if err != nil {
return fmt.Errorf("failed to delete user tokens: %w", err)
}

// Log the action
s.authLogRepo.Create(&models.AuthLog{
UserID:    userID,
Action:    models.AuthActionLogoutAll,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
})

return nil
}

// ChangePassword updates a user's password
func (s *AuthService) ChangePassword(userID, currentPassword, newPassword string, ipAddress, userAgent string) error {
// Get user
user, err := s.userRepo.GetByID(userID)
if err != nil {
return models.ErrUserNotFound
}

// Verify current password
err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
if err != nil {
// Log failed password change attempt
s.authLogRepo.Create(&models.AuthLog{
	UserID:    userID,
	Action:    models.AuthActionChangePassword,
	IPAddress: ipAddress,
	UserAgent: userAgent,
	Status:    models.StatusFailure,
	Details:   "Current password verification failed",
})
return models.ErrInvalidCredentials
}

// Hash new password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
if err != nil {
return fmt.Errorf("failed to hash password: %w", err)
}

// Update password in database
user.PasswordHash = string(hashedPassword)
err = s.userRepo.UpdatePassword(userID, string(hashedPassword))
if err != nil {
return fmt.Errorf("failed to update password: %w", err)
}

// Log successful password change
s.authLogRepo.Create(&models.AuthLog{
UserID:    userID,
Action:    models.AuthActionChangePassword,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
})

// Invalidate all existing sessions to force re-login with new password
return s.LogoutAllSessions(userID, ipAddress, userAgent)
}

// RequestPasswordReset initiates a password reset process
func (s *AuthService) RequestPasswordReset(email, ipAddress, userAgent string) (string, error) {
// Get user by email
user, err := s.userRepo.GetByEmail(email)
if err != nil {
// We don't want to leak information about whether an email exists
// But we still log the attempt
s.authLogRepo.Create(&models.AuthLog{
	Action:    models.AuthActionPasswordResetRequest,
	IPAddress: ipAddress,
	UserAgent: userAgent,
	Status:    models.StatusFailure,
	Details:   "Email not found",
})
return "", nil
}

// Generate reset token
tokenBytes := make([]byte, 32)
_, err = rand.Read(tokenBytes)
if err != nil {
return "", fmt.Errorf("failed to generate reset token: %w", err)
}
resetToken := base64.URLEncoding.EncodeToString(tokenBytes)

// Hash token for storage
hashedToken, err := bcrypt.GenerateFromPassword([]byte(resetToken), bcrypt.DefaultCost)
if err != nil {
return "", fmt.Errorf("failed to hash reset token: %w", err)
}

// Store reset token with expiry (usually 24 hours)
expiry := time.Now().Add(24 * time.Hour)
err = s.userRepo.StorePasswordResetToken(user.ID, string(hashedToken), expiry)
if err != nil {
return "", fmt.Errorf("failed to store reset token: %w", err)
}

// Log the request
s.authLogRepo.Create(&models.AuthLog{
UserID:    user.ID,
Action:    models.AuthActionPasswordResetRequest,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
})

return resetToken, nil
}

// ResetPassword completes the password reset process
func (s *AuthService) ResetPassword(email, resetToken, newPassword, ipAddress, userAgent string) error {
// Get user by email
user, err := s.userRepo.GetByEmail(email)
if err != nil {
return models.ErrUserNotFound
}

// Get stored reset token
storedToken, expiry, err := s.userRepo.GetPasswordResetToken(user.ID)
if err != nil {
return models.ErrInvalidToken
}

// Check if token is expired
if expiry.Before(time.Now()) {
s.authLogRepo.Create(&models.AuthLog{
	UserID:    user.ID,
	Action:    models.AuthActionPasswordReset,
	IPAddress: ipAddress,
	UserAgent: userAgent,
	Status:    models.StatusFailure,
	Details:   "Reset token expired",
})
return models.ErrTokenExpired
}

// Verify token
err = bcrypt.CompareHashAndPassword([]byte(storedToken), []byte(resetToken))
if err != nil {
s.authLogRepo.Create(&models.AuthLog{
	UserID:    user.ID,
	Action:    models.AuthActionPasswordReset,
	IPAddress: ipAddress,
	UserAgent: userAgent,
	Status:    models.StatusFailure,
	Details:   "Invalid reset token",
})
return models.ErrInvalidToken
}

// Hash new password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
if err != nil {
return fmt.Errorf("failed to hash password: %w", err)
}

// Update password
err = s.userRepo.UpdatePassword(user.ID, string(hashedPassword))
if err != nil {
return fmt.Errorf("failed to update password: %w", err)
}

// Clear reset token
err = s.userRepo.ClearPasswordResetToken(user.ID)
if err != nil {
return fmt.Errorf("failed to clear reset token: %w", err)
}

// Invalidate all existing sessions
err = s.tokenRepo.DeleteAllForUser(user.ID)
if err != nil {
return fmt.Errorf("failed to invalidate sessions: %w", err)
}

// Log successful password reset
s.authLogRepo.Create(&models.AuthLog{
UserID:    user.ID,
Action:    models.AuthActionPasswordReset,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
})

return nil
}

// ToggleTwoFactor enables or disables 2FA for a user
func (s *AuthService) ToggleTwoFactor(userID string, enable bool, ipAddress, userAgent string) error {
// Update user's 2FA status
err := s.userRepo.UpdateTwoFactorStatus(userID, enable)
if err != nil {
return fmt.Errorf("failed to update 2FA status: %w", err)
}

// Log the action
action := models.AuthActionEnable2FA
if !enable {
action = models.AuthActionDisable2FA
}

s.authLogRepo.Create(&models.AuthLog{
UserID:    userID,
Action:    action,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
})

return nil
}

// ValidateToken checks if a token is valid
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
// Parse and validate token
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
}
return []byte(s.config.JWTSecret), nil
})

if err != nil || !token.Valid {
return nil, models.ErrInvalidToken
}

claims, ok := token.Claims.(jwt.MapClaims)
if !ok {
return nil, models.ErrInvalidToken
}

// Check token type
tokenType, ok := claims["type"].(string)
if !ok || tokenType != "access" {
return nil, models.ErrInvalidToken
}

// Get user ID from token
userID, ok := claims["sub"].(string)
if !ok {
return nil, models.ErrInvalidToken
}

// Check if token exists in database
tokenHash, ok := claims["jti"].(string)
if !ok {
return nil, models.ErrInvalidToken
}

storedToken, err := s.tokenRepo.GetByHash(tokenHash)
if err != nil || storedToken == nil {
return nil, models.ErrInvalidToken
}

// Check if token is expired (in database)
if storedToken.ExpiresAt.Before(time.Now()) {
return nil, models.ErrTokenExpired
}

// Get user
user, err := s.userRepo.GetByID(userID)
if err != nil {
return nil, models.ErrUserNotFound
}

// Check account status
if user.AccountStatus != models.UserStatusActive {
return nil, models.ErrAccountNotActive
}

return user, nil
}

// UpdateUserRole updates a user's role
func (s *AuthService) UpdateUserRole(userID string, newRole models.UserRole, adminID, ipAddress, userAgent string) error {
// Check if admin exists
admin, err := s.userRepo.GetByID(adminID)
if err != nil {
return models.ErrUserNotFound
}

// Check admin permissions
if admin.Role != models.UserRoleAdmin && admin.Role != models.UserRoleSuperAdmin {
return models.ErrUnauthorized
}

// Get user
user, err := s.userRepo.GetByID(userID)
if err != nil {
return models.ErrUserNotFound
}

// Prevent role escalation to admin by non-superadmins
if (newRole == models.UserRoleAdmin || newRole == models.UserRoleSuperAdmin) && admin.Role != models.UserRoleSuperAdmin {
return models.ErrUnauthorized
}

// Update user role
user.Role = newRole
err = s.userRepo.UpdateRole(userID, newRole)
if err != nil {
return fmt.Errorf("failed to update role: %w", err)
}

// Log the change
s.authLogRepo.Create(&models.AuthLog{
UserID:    adminID,
Action:    models.AuthActionUpdateRole,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
Details:   fmt.Sprintf("Updated user %s role to %s", userID, newRole),
})

return nil
}

// UpdateAccountStatus changes a user's account status
func (s *AuthService) UpdateAccountStatus(userID string, status models.UserStatus, adminID, ipAddress, userAgent string) error {
// Check if admin exists
admin, err := s.userRepo.GetByID(adminID)
if err != nil {
return models.ErrUserNotFound
}

// Check admin permissions
if admin.Role != models.UserRoleAdmin && admin.Role != models.UserRoleSuperAdmin {
return models.ErrUnauthorized
}

// Get user
user, err := s.userRepo.GetByID(userID)
if err != nil {
return models.ErrUserNotFound
}

// Prevent modifying super admin accounts by non-super admins
if user.Role == models.UserRoleSuperAdmin && admin.Role != models.UserRoleSuperAdmin {
return models.ErrUnauthorized
}

// Update account status
user.AccountStatus = status
err = s.userRepo.UpdateStatus(userID, status)
if err != nil {
return fmt.Errorf("failed to update account status: %w", err)
}

// Log the change
s.authLogRepo.Create(&models.AuthLog{
UserID:    adminID,
Action:    models.AuthActionUpdateStatus,
IPAddress: ipAddress,
UserAgent: userAgent,
Status:    models.StatusSuccess,
Details:   fmt.Sprintf("Updated user %s status to %s", userID, status),
})

// If account is suspended or deactivated, invalidate all sessions
if status == models.UserStatusSuspended || status == models.UserStatusDeactivated {
s.tokenRepo.DeleteAllForUser(userID)
}

return nil
}

// generateTokenPair creates a new pair of access and refresh tokens
func (s *AuthService) generateTokenPair(user *models.User) (*models.TokenPair, error) {
// Generate unique token IDs
accessJTI, err := generateRandomString(32)
if err != nil {
return nil, err
}

refreshJTI, err := generateRandomString(32)
if err != nil {
return nil, err
}

// Generate access token
accessTokenClaims := jwt.MapClaims{
"sub": user.ID,
"exp": time.Now().Add(time.Duration(s.config.AccessTokenLifetimeMinutes) * time.Minute).Unix(),
"iat": time.Now().Unix(),
"jti": accessJTI,
"type": "access",
"role": user.Role,
"email": user.Email,
"name": user.Name,
}

accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
if err != nil {
return nil, fmt.Errorf("failed to sign access token: %w", err)
}

// Generate refresh token
refreshTokenClaims := jwt.MapClaims{
"sub": user.ID,
"exp": time.Now().Add(time.Duration(s.config.RefreshTokenLifetimeDays) * 24 * time.Hour).Unix(),
"iat": time.Now().Unix(),
"jti": refreshJTI,
"type": "refresh",
}

refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
if err != nil {
return nil, fmt.Errorf("failed to sign refresh token: %w", err)
}

return &models.TokenPair{
AccessToken:       accessTokenString,
RefreshToken:      refreshTokenString,
AccessTokenHash:   accessJTI,
RefreshTokenHash:  refreshJTI,
}, nil
}

// generateRandomString creates a random string of specified length
func generateRandomString(length int) (string, error) {
bytes := make([]byte, length)
if _, err := rand.Read(bytes); err != nil {
return "", err
}
return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}