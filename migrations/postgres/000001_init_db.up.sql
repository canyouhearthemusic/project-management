CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	registration_date DATE NOT NULL,
	role VARCHAR CHECK (role IN ('admin', 'manager', 'developer'))
);

CREATE INDEX IF NOT EXISTS users_name_idx ON users(email);

CREATE TABLE IF NOT EXISTS projects (
	id VARCHAR(255) PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL,
	started_at DATE NOT NULL,
	finished_at DATE NOT NULL,
	manager_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS projects_title_idx ON projects(title);
CREATE INDEX IF NOT EXISTS projects_manager_idx ON projects(manager_id);

CREATE TABLE IF NOT EXISTS tasks (
	id VARCHAR(255) PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL,
	priority VARCHAR CHECK (priority IN ('low', 'medium', 'high')),
	status VARCHAR CHECK (status IN ('active', 'in_proccess', 'done')),
	author_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
	project_id VARCHAR(255) REFERENCES projects(id) ON DELETE CASCADE,
	created_at DATE NOT NULL,
	done_at DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS tasks_title_idx ON tasks(title);
CREATE INDEX IF NOT EXISTS tasks_author_idx ON tasks(author_id);
CREATE INDEX IF NOT EXISTS tasks_project_idx ON tasks(project_id);
CREATE INDEX IF NOT EXISTS tasks_status_idx ON tasks(status);
CREATE INDEX IF NOT EXISTS tasks_priority_idx ON tasks(priority);
