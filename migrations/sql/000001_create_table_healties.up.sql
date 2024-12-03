CREATE TABLE IF NOT EXISTS healties(
    id BIGINT AUTO_INCREMENT,
    name VARCHAR(50) NULL,
    status BOOLEAN,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT NULL,
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_by BIGINT NULL,
    CONSTRAINT healty_pk PRIMARY KEY(id)
);

INSERT INTO healties(
    id,
    name,
    status,
    created_at,
    updated_at,
    created_by,
    updated_by,
    deleted_by
)
VALUES(1, 'test pertama', TRUE, 1, 1, 1, 1, 0),(2, 'test kedua', FALSE, 1, 1, 1, 1, 0);