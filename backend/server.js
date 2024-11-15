import express from 'express';
import dotenv from 'dotenv';
import mysql from 'mysql2/promise';
import serverless from 'serverless-http';
import { generate } from 'random-words';

dotenv.config();

const app = express();

// Middleware to parse JSON requests
app.use(express.json());

// Function to generate a unique game ID from 4 random words
function generateGameID() {
  const words = generate({ exactly: 4, maxLength: 5 }); // Generate 4 random words
  return words.join('-'); // Join them with hyphens
}

// Define migrations
const migrations = [
  {
    version: 1,
    sql: `
      CREATE TABLE IF NOT EXISTS games (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        timeCreated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );

      CREATE TABLE IF NOT EXISTS players (
        gameId VARCHAR(255),
        playerName VARCHAR(255) NOT NULL UNIQUE,
        team INT,
        timeJoined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (gameId) REFERENCES games(id) ON DELETE CASCADE
      );

      CREATE TABLE IF NOT EXISTS cards (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        title VARCHAR(255) NOT NULL UNIQUE,
        description TEXT,
        playerNameCreatedBy VARCHAR(255),
        FOREIGN KEY (playerNameCreatedBy) REFERENCES players(playerName) ON DELETE SET NULL
      );

      CREATE TABLE IF NOT EXISTS turns (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        gameId VARCHAR(255),
        roundNumber INT,
        turnNumber INT,
        timeStarted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        playerName VARCHAR(255),
        FOREIGN KEY (gameId) REFERENCES games(id) ON DELETE CASCADE,
        FOREIGN KEY (playerName) REFERENCES players(playerName) ON DELETE SET NULL
      );

      CREATE TABLE IF NOT EXISTS wonCards (
        turnId VARCHAR(255),
        cardId VARCHAR(255),
        FOREIGN KEY (turnId) REFERENCES turns(id) ON DELETE CASCADE,
        FOREIGN KEY (cardId) REFERENCES cards(id) ON DELETE CASCADE,
        PRIMARY KEY (turnId, cardId)
      );

      CREATE TABLE IF NOT EXISTS skippedCards (
        turnId VARCHAR(255),
        cardId VARCHAR(255),
        FOREIGN KEY (turnId) REFERENCES turns(id) ON DELETE CASCADE,
        FOREIGN KEY (cardId) REFERENCES cards(id) ON DELETE CASCADE,
        PRIMARY KEY (turnId, cardId)
      );
    `
  }
  // Add new migrations here
];

// Function to apply migrations
async function applyMigrations(connection) {
  const migrationsTable = 'schema_migrations';

  // Check if the database exists
  const [databases] = await connection.query('SHOW DATABASES LIKE ?', [process.env.DB_NAME]);
  if (databases.length === 0) {
    // Create the database if it doesn't exist
    await connection.query(`CREATE DATABASE \`${process.env.DB_NAME}\``);
    console.log(`Database ${process.env.DB_NAME} created.`);
  }

  // Connect to the specific database
  await connection.changeUser({ database: process.env.DB_NAME });

  // Ensure schema_migrations table exists
  await connection.query(`
    CREATE TABLE IF NOT EXISTS ${migrationsTable} (
      version INT NOT NULL,
      applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
  `);

  for (const migration of migrations) {
    const [rows] = await connection.query(`SELECT version FROM ${migrationsTable} WHERE version = ?`, [migration.version]);
    if (rows.length === 0) {
      // Split the SQL statements and execute them one by one
      const sqlStatements = migration.sql.split(';').filter(stmt => stmt.trim());
      for (const statement of sqlStatements) {
        await connection.query(statement);
      }
      // Record migration
      await connection.query(`INSERT INTO ${migrationsTable} (version) VALUES (?)`, [migration.version]);
      console.log(`Migration ${migration.version} applied successfully`);
    }
  }
}

// Define a route for creating a new game using GET
app.get('/new-game', async (req, res) => {
  const playerName = req.query.playerName; // Retrieve playerName from query parameters
  const gameID = generateGameID();
  const now = new Date();

  try {
    // Connect to MySQL server without specifying a database
    const connection = await mysql.createConnection({
      host: process.env.DB_HOST,
      user: process.env.DB_USER,
      password: process.env.DB_PASSWORD
    });

    await applyMigrations(connection);

    const gameInsertSQL = 'INSERT INTO games (id, timeCreated) VALUES (?, ?)';
    await connection.execute(gameInsertSQL, [gameID, now]);

    const playerInsertSQL = 'INSERT INTO players (gameId, playerName, team, timeJoined) VALUES (?, ?, ?, ?)';
    await connection.execute(playerInsertSQL, [gameID, playerName, 0, now]);

    await connection.end();

    res.json({ gameId: gameID });
  } catch (err) {
    console.error('Database operation failed:', err);
    res.status(500).send('Failed to create game');
  }
});

// Export the app wrapped in serverless-http
export const handler = serverless(app); 