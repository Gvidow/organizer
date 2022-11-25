CREATE TABLE user(
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    user_login VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE task(
    task_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    task_date DATE NOT NULL,
    task_position INT NOT NULL,
    task_title VARCHAR(255),
    task_description VARCHAR(255),
    FOREIGN KEY (user_id)
        REFERENCES user(user_id)
        ON DELETE CASCADE
);

CREATE TABLE _session(
    session_id VARCHAR(65) PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    use_date DATE NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES user(user_id)
        ON DELETE CASCADE
);
