import express from 'express';
import dotenv from 'dotenv';
import mysql from 'mysql2/promise';
import serverless from 'serverless-http';
import { generate } from 'random-words';

dotenv.config();

const app = express();

// Middleware to parse JSON requests
app.use(express.json());

// // Define the new GameState type
// interface GameState {
//   gameId: string;
//   teams: { players: string[]; score: number }[];
//   cards: { title: string; description: string; cardId: number }[];
//   timeLeft: number;
//   playerTurn: string;
//   state: 'JOIN' | 'READY' | 'ROUND' | 'COMPLETE'; // You can add more states if needed
//   roundNumber: number;
// }

// Function to generate a unique game ID from 4 random words
function generateGameID() { // Added return type
  const words = generate({ exactly: 4, maxLength: 5 }); // Generate 4 random words
  return words.join('-'); // Join them with hyphens
}

// Define migrations
const migrations = [ // Added type for migrations
  {
    version: 1,
    sql: `
      CREATE TABLE IF NOT EXISTS games (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        timeCreated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );

      CREATE TABLE IF NOT EXISTS players (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        gameId VARCHAR(255),
        playerName VARCHAR(255) NOT NULL,
        team INT,
        timeJoined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (gameId) REFERENCES games(id) ON DELETE CASCADE,
        UNIQUE (gameId, playerName)
      );

      CREATE TABLE IF NOT EXISTS cards (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT,
        playerIdCreatedBy INT,
        FOREIGN KEY (playerIdCreatedBy) REFERENCES players(id) ON DELETE SET NULL
      );

      CREATE TABLE IF NOT EXISTS turns (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        gameId VARCHAR(255),
        roundNumber INT,
        turnNumber INT,
        timeStarted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        playerId INT,
        FOREIGN KEY (gameId) REFERENCES games(id) ON DELETE CASCADE,
        FOREIGN KEY (playerId) REFERENCES players(id) ON DELETE SET NULL
      );

      CREATE TABLE IF NOT EXISTS wonCards (
        turnId INT,
        cardId INT,
        FOREIGN KEY (turnId) REFERENCES turns(id) ON DELETE CASCADE,
        FOREIGN KEY (cardId) REFERENCES cards(id) ON DELETE CASCADE,
        PRIMARY KEY (turnId, cardId)
      );

      CREATE TABLE IF NOT EXISTS skippedCards (
        turnId INT,
        cardId INT,
        FOREIGN KEY (turnId) REFERENCES turns(id) ON DELETE CASCADE,
        FOREIGN KEY (cardId) REFERENCES cards(id) ON DELETE CASCADE,
        PRIMARY KEY (turnId, cardId)
      );
    `
  }
  // Add new migrations here
];

// Function to apply migrations
async function applyMigrations(connection) { // Added types
  const migrationsTable = 'schema_migrations';

  // Temporary - delete old database
  await connection.query(`DROP DATABASE IF EXISTS \`${process.env.DB_NAME}\``);
  console.log(`Database ${process.env.DB_NAME} deleted.`);

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

  // Get the most recent migration version
  const [recentVersionRows] = await connection.query(`SELECT MAX(version) AS version FROM ${migrationsTable}`);
  const recentVersion = recentVersionRows[0].version || 0; // Default to 0 if no migrations exist

  // Apply new migrations
  for (const migration of migrations) {
    if (migration.version > recentVersion) {
      // Split the SQL statements and execute them one by one
      const sqlStatements = migration.sql.split(';').filter(sql => sql.trim() !== ''); // Filter out empty statements
      for (const sql of sqlStatements) {
        await connection.query(sql);
      }
      await connection.query(`INSERT INTO ${migrationsTable} (version) VALUES (?)`, [migration.version]);
      console.log(`Applied migration version ${migration.version}`);
    }
  }
}

// Define a route for creating a new game using GET
app.get('/new-game', async (req, res) => { // Added types for req and res
  const playerName = req.query.playerName; // Added type assertion
  const gameID = generateGameID(); // Added type
  const now = new Date(); // Added type

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