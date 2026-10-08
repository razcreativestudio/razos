# Skema Database Internal RAZOS

RAZOS menggunakan SQLite murni yang disimpan di `/DATA/razos/data/razos.db` untuk mencatat metadata sistem tanpa daemon tambahan.

```sql
-- Tabel Kredensial & Pengguna
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Pengaturan Sistem
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Tabel Manajemen Website
CREATE TABLE IF NOT EXISTS websites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    domain TEXT UNIQUE NOT NULL,
    container_id TEXT NOT NULL,
    host_port INTEGER UNIQUE NOT NULL,
    runtime TEXT NOT NULL, -- 'php', 'nodejs', 'static'
    database_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(database_id) REFERENCES databases(id) ON DELETE SET NULL
);

-- Tabel Database Terdaftar
CREATE TABLE IF NOT EXISTS databases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    engine TEXT NOT NULL, -- 'mariadb', 'postgres'
    db_user TEXT NOT NULL,
    db_password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel App Store Terpasang
CREATE TABLE IF NOT EXISTS installed_apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_slug TEXT UNIQUE NOT NULL,
    app_name TEXT NOT NULL,
    container_name TEXT NOT NULL,
    web_port INTEGER,
    data_path TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```
