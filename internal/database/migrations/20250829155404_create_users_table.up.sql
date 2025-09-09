CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');

CREATE TABLE IF NOT EXISTS users(
   id bigserial PRIMARY KEY,
   first_name VARCHAR (50) NOT NULL,
   last_name VARCHAR (50) NOT NULL,
   password VARCHAR (255) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   email_verified_at TIMESTAMP WITH TIME ZONE,
   status user_status NOT NULL DEFAULT 'active',
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP WITH TIME ZONE
);

-- User table indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Auto Update updated_at on row update
CREATE TRIGGER users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Soft delete implementation
CREATE TRIGGER users_soft_delete
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION soft_delete();

-- Failed Logins Table
CREATE TABLE IF NOT EXISTS failed_logins(
   id bigserial PRIMARY KEY,
   user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
   ip_address INET NOT NULL,
   attempted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Auth Tokens Table
CREATE TABLE IF NOT EXISTS auth_tokens(
   id bigserial PRIMARY KEY,
   user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
   type VARCHAR(20),
   token VARCHAR(255),
   expire_at TIMESTAMP WITH TIME ZONE,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

