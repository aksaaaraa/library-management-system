# frontend

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Run your unit tests
```
npm run test:unit
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).


```
library-management-system/
├── backend/                 # Golang backend
│   ├── cmd/                 # Entry points
│   │   └── server/          # API server
│   ├── internal/            # Private application code
│   │   ├── api/             # API handlers
│   │   ├── auth/            # Authentication logic
│   │   ├── config/          # Configuration
│   │   ├── middleware/      # Middleware components
│   │   ├── models/          # Data models
│   │   ├── repository/      # Database repositories
│   │   ├── service/         # Business logic
│   │   └── validation/      # Input validation
│   ├── migrations/          # Database migrations
│   ├── pkg/                 # Public libraries
│   │   ├── logger/          # Logging utilities
│   │   └── security/        # Security utilities
│   ├── scripts/             # Utility scripts
│   ├── .env.example         # Example environment variables
│   ├── .gitignore           # Git ignore file
│   ├── go.mod               # Go modules file
│   ├── go.sum               # Go modules checksum
│   └── main.go              # Main application entry
├── frontend/                # Vue.js frontend
│   ├── public/              # Static files
│   ├── src/                 # Source files
│   │   ├── assets/          # Assets (images, fonts)
│   │   ├── components/      # Vue components
│   │   ├── layouts/         # Layout components
│   │   ├── router/          # Vue Router configuration
│   │   ├── services/        # API services
│   │   ├── store/           # Vuex store
│   │   ├── utils/           # Utility functions
│   │   ├── views/           # Vue views/pages
│   │   ├── App.vue          # Root component
│   │   └── main.js          # Entry point
│   ├── .gitignore           # Git ignore file
│   ├── package.json         # NPM package file
│   └── tailwind.config.js   # Tailwind CSS configuration
├── docs/                    # Documentation
│   └── swagger/             # Swagger API documentation
├── .github/                 # GitHub workflows
│   └── workflows/           # CI/CD workflows
├── .gitignore               # Git ignore file
└── README.md                # Project documentation
```

-- Database schema for the Library Management System

-- Users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'librarian', 'member') NOT NULL DEFAULT 'member',
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    account_status ENUM('active', 'suspended', 'inactive') DEFAULT 'active',
    failed_login_attempts INT DEFAULT 0,
    reset_token VARCHAR(255) NULL,
    reset_token_expires TIMESTAMP NULL,
    two_factor_secret VARCHAR(255) NULL,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_status (account_status)
) ENGINE=InnoDB;

-- Books table
CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    isbn VARCHAR(20) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(255),
    publication_year INT,
    description TEXT,
    category VARCHAR(100),
    language VARCHAR(50),
    page_count INT,
    total_copies INT NOT NULL DEFAULT 0,
    available_copies INT NOT NULL DEFAULT 0,
    location VARCHAR(100),
    cover_image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT,
    INDEX idx_isbn (isbn),
    INDEX idx_title (title),
    INDEX idx_author (author),
    INDEX idx_category (category),
    INDEX idx_availability (available_copies),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- Book copies table (for tracking individual book copies)
CREATE TABLE book_copies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    book_id INT NOT NULL,
    copy_number VARCHAR(50) NOT NULL,
    status ENUM('available', 'borrowed', 'reserved', 'lost', 'damaged', 'in_repair') DEFAULT 'available',
    acquisition_date DATE,
    condition ENUM('new', 'good', 'fair', 'poor') DEFAULT 'new',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (book_id, copy_number),
    INDEX idx_status (status),
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Borrowing records table
CREATE TABLE borrowings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    book_copy_id INT NOT NULL,
    borrowed_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP NOT NULL,
    returned_date TIMESTAMP NULL,
    status ENUM('active', 'returned', 'overdue', 'lost') DEFAULT 'active',
    fine_amount DECIMAL(10, 2) DEFAULT 0.00,
    fine_paid BOOLEAN DEFAULT FALSE,
    staff_id_checkout INT,
    staff_id_return INT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_book_copy_id (book_copy_id),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_copy_id) REFERENCES book_copies(id) ON DELETE CASCADE,
    FOREIGN KEY (staff_id_checkout) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (staff_id_return) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- Reservations table
CREATE TABLE reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    reservation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expiry_date TIMESTAMP NOT NULL,
    status ENUM('pending', 'fulfilled', 'cancelled', 'expired') DEFAULT 'pending',
    fulfilled_date TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_book_id (book_id),
    INDEX idx_status (status),
    INDEX idx_expiry_date (expiry_date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Authentication logs
CREATE TABLE auth_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    action ENUM('login', 'logout', 'failed_login', 'password_reset', 'password_change', '2fa_setup', '2fa_verify') NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    status ENUM('success', 'failure') NOT NULL,
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- API access tokens
CREATE TABLE access_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    token_type ENUM('access', 'refresh') NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Activity logs for auditing
CREATE TABLE activity_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT,
    old_value TEXT,
    new_value TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_entity_type (entity_type),
    INDEX idx_entity_id (entity_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- System settings
CREATE TABLE settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by INT,
    INDEX idx_setting_key (setting_key),
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;