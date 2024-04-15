CREATE TABLE user (
    user_full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    PRIMARY KEY (user_full_name)
);

CREATE TABLE customer (
    customer_id INT NOT NULL AUTO_INCREMENT,
    account_number VARCHAR(20) NOT NULL,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL,
    middle_name VARCHAR(20) NOT NULL,
    date_of_birth DATE NOT NULL,
    tax_id VARCHAR(12) NOT NULL,
    user_full_name VARCHAR(120) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Не в работе',
    PRIMARY KEY (customer_id),
    FOREIGN KEY (user_full_name) REFERENCES user(user_full_name)
);