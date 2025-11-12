-- This SQL script creates a table named 'otps' in the database.
-- The table is designed to store OTP (One-Time Password) information.
-- Indexes are added for efficient lookup during validation.
CREATE TABLE IF NOT EXISTS otps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,           -- Auto-incrementing ID
    user_id VARCHAR(50) NOT NULL,                   -- Reference to the user (short identifier)
    otp_code CHAR(6) NOT NULL,                      -- OTP code (6 digits)
    status TINYINT DEFAULT 1,                       -- OTP status (e.g. 1 = created, 2 = validated, 3 = expired), see the application code.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically set creation timestamp
    expires_at TIMESTAMP NOT NULL,                  -- OTP expiration timestamp
    validated_at TIMESTAMP NULL,                    -- When OTP was successfully validated

    CONSTRAINT uq_otp_user_code UNIQUE(user_id, otp_code)  -- Prevent duplicate OTPs for the same user
);
