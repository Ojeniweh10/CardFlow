-- ============================================================
-- Extensions
-- ============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================================
-- Users
-- ============================================================

CREATE TABLE users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email               VARCHAR(255) NOT NULL UNIQUE,
    password_hash       VARCHAR(255) NOT NULL,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    phone               VARCHAR(20),
    status              VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'closed')),
    email_verified      BOOLEAN NOT NULL DEFAULT FALSE,
    otp                 VARCHAR(6),
    otp_expires_at      TIMESTAMP,
    mfa_enabled         BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_secret          VARCHAR(255),
    last_login_at       TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email        ON users(email);
CREATE INDEX idx_users_status       ON users(status);
CREATE INDEX idx_users_created_at   ON users(created_at);

-- ============================================================
-- Admins
-- ============================================================

CREATE TABLE admins (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email               VARCHAR(255) NOT NULL UNIQUE,
    password_hash       VARCHAR(255) NOT NULL,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    phone               VARCHAR(20),
    role                VARCHAR(50) NOT NULL DEFAULT 'admin' CHECK (role IN ('superadmin', 'admin', 'compliance_officer')),
    status              VARCHAR(50) NOT NULL DEFAULT 'active'CHECK (status IN ('active', 'inactive')),
    email_verified      BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_enabled         BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_secret          VARCHAR(255),
    last_login_at       TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admins_email       ON admins(email);
CREATE INDEX idx_admins_status      ON admins(status);
CREATE INDEX idx_admins_created_at  ON admins(created_at);

-- ============================================================
-- KYC
-- ============================================================

CREATE TABLE kyc_submissions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status              VARCHAR(50) NOT NULL CHECK (status IN ('started','documents_uploaded', 'under_review', 'verified', 'rejected', 'expired')),
    rejection_reason    TEXT,
    reviewed_by         UUID REFERENCES admins(id),
    submitted_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at         TIMESTAMP,
    expires_at          TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kyc_user_id        ON kyc_submissions(user_id);
CREATE INDEX idx_kyc_status         ON kyc_submissions(status);

CREATE TABLE kyc_documents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kyc_submission_id   UUID NOT NULL UNIQUE REFERENCES kyc_submissions(id) ON DELETE CASCADE,
    document_type       VARCHAR(50) NOT NULL UNIQUE CHECK (document_type IN ('id_document', 'proof_of_address', 'selfie')),
    mime_type           VARCHAR(100) NOT NULL,
    encrypted_data      BYTEA NOT NULL,
    encryption_version  VARCHAR(20) NOT NULL DEFAULT 'v1',
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- Cards
-- ============================================================

CREATE TABLE cards (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_reference          VARCHAR(100) NOT NULL UNIQUE,
    masked_pan              VARCHAR(19) NOT NULL,
    last_four               VARCHAR(4) NOT NULL,
    pan_hash                VARCHAR(255) NOT NULL,
    cvv_hash                VARCHAR(255) NOT NULL,
    card_type               VARCHAR(50) CHECK (card_type IN ('single-use', 'multi-use')),
    currency                VARCHAR(3) NOT NULL DEFAULT 'USD',
    status                  VARCHAR(50) CHECK (status IN ('active', 'frozen', 'terminated', 'expired')),
    spending_limit_amount   DECIMAL(15,2),
    spending_limit_period   VARCHAR(20),
    current_balance         DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    expiry_month            VARCHAR(2),
    expiry_year             VARCHAR(4),
    expires_at              TIMESTAMP,
    issued_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cards_user_id      ON cards(user_id);
CREATE INDEX idx_cards_status       ON cards(status);

-- ============================================================
-- Transactions
-- ============================================================

CREATE TABLE transactions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id                 UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_reference   VARCHAR(100) NOT NULL UNIQUE,
    idempotency_key         VARCHAR(100),
    amount                  DECIMAL(15,2) NOT NULL,
    authorized_amount       DECIMAL(15,2),
    captured_amount         DECIMAL(15,2),
    currency                VARCHAR(3) NOT NULL,
    merchant_name           VARCHAR(255),
    merchant_mcc            VARCHAR(4),
    merchant_country        VARCHAR(2),
    status                  VARCHAR(50),
    type                    VARCHAR(50),
    decline_reason          TEXT,
    metadata_json           JSONB,
    transaction_timestamp   TIMESTAMP NOT NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_card_id   ON transactions(card_id);
CREATE INDEX idx_transactions_user_id   ON transactions(user_id);
CREATE INDEX idx_transactions_created   ON transactions(created_at);

-- ============================================================
-- Balance Ledger (Source of Truth)
-- ============================================================

CREATE TABLE balance_ledger (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id         UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    transaction_id  UUID REFERENCES transactions(id),
    entry_type      VARCHAR(50),
    amount          DECIMAL(15,2) NOT NULL,
    balance_after   DECIMAL(15,2) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ledger_card_id      ON balance_ledger(card_id);
CREATE INDEX idx_ledger_created_at   ON balance_ledger(created_at);

-- ============================================================
-- Audit Logs
-- ============================================================

CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    action          VARCHAR(100) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       UUID,
    ip_address      INET,
    user_agent      TEXT,
    request_id      VARCHAR(100),
    metadata_json   JSONB,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user_id       ON audit_logs(user_id);
CREATE INDEX idx_audit_action        ON audit_logs(action);
CREATE INDEX idx_audit_entity        ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_created_at    ON audit_logs(created_at);

-- ============================================================
-- Notifications
-- ============================================================

CREATE TABLE notifications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL CHECK (type IN ('transaction', 'kyc', 'card', 'security', 'system')),
    channel         VARCHAR(20) NOT NULL CHECK (channel IN ('email', 'sms', 'push')),
    subject         VARCHAR(255),
    body            TEXT NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed')),
    retry_count     INT NOT NULL DEFAULT 0,
    error_message   TEXT,
    sent_at         TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id     ON notifications(user_id);
CREATE INDEX idx_notifications_status      ON notifications(status);
CREATE INDEX idx_notifications_created_at  ON notifications(created_at);

-- ============================================================
-- Refresh Tokens
-- ============================================================

CREATE TABLE refresh_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) NOT NULL UNIQUE,
    expires_at      TIMESTAMP NOT NULL,
    revoked         BOOLEAN NOT NULL DEFAULT FALSE,
    revoked_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refresh_user_id     ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_token_hash  ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_expires_at  ON refresh_tokens(expires_at);
