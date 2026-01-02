-- Factorial jobs metadata
CREATE TABLE factorial_jobs (
    job_id      VARCHAR(36) PRIMARY KEY,
    total_items INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Results per item
CREATE TABLE factorial_results (
    job_id     VARCHAR(36) NOT NULL REFERENCES factorial_jobs(job_id) ON DELETE CASCADE,
    item_id    INTEGER NOT NULL,
    input      INTEGER NOT NULL,
    output     NUMERIC(1000) NULL,
    status     VARCHAR(8) NOT NULL CHECK (status IN ('done','failed')),
    error      VARCHAR(30) NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (job_id, item_id)
);

CREATE INDEX idx_factorial_results_job_id ON factorial_results(job_id);
CREATE INDEX idx_factorial_results_job_status ON factorial_results(job_id, status);

