ALTER TABLE users ADD COLUMN fcm_token VARCHAR(255);

CREATE INDEX idx_users_fcm_token ON users(fcm_token);