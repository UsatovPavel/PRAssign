-- Index to speed up retention queries by updated_at
CREATE INDEX idx_factorial_results_updated_at ON factorial_results(updated_at);

