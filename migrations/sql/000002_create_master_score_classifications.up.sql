CREATE TABLE IF NOT EXISTS master_score_classifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    min_score INT NOT NULL,
    max_score INT DEFAULT NULL,
    code VARCHAR(50) NOT NULL UNIQUE, -- e.g. "BURUK", "BUTUH_PENGEMBANGAN", etc
    label VARCHAR(100) NOT NULL,      -- e.g. "Buruk"
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

INSERT INTO master_score_classifications (min_score, max_score, code, label, description)
VALUES
  (NULL, 49, 'BURUK', 'Buruk', 'Nilai kurang dari 50'),
  (50, 59, 'BUTUH_PENGEMBANGAN', 'Butuh Pengembangan', 'Nilai antara 50 - 59'),
  (60, 69, 'CUKUP', 'Cukup', 'Nilai antara 60 - 69'),
  (70, 79, 'BAIK', 'Baik', 'Nilai antara 70 - 79'),
  (80, 89, 'SANGAT_BAIK', 'Sangat Baik', 'Nilai antara 80 - 89'),
  (90, NULL, 'ISTIMEWA', 'Istimewa', 'Nilai >= 90');
