Set up MongoDB and Mongo Express using Docker Compose for the Vajra project.

1. Docker Compose configuration

Use the following configuration format and requirements:

services:
  mongo:
    image: mongo:latest
    container_name: vajra-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: vijay_mehra
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD}
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    networks:
      - mongo-db-network

  mongo-express:
    image: mongo-express:latest
    container_name: vajra-mongo-express
    restart: always
    depends_on:
      - mongo
    ports:
      - "8081:8081"
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: vijay_mehra
      ME_CONFIG_MONGODB_ADMINPASSWORD: ${MONGO_PASSWORD}
      ME_CONFIG_MONGODB_PORT: 27017
      ME_CONFIG_MONGODB_SERVER: mongo
      ME_CONFIG_MONGODB_AUTH_DATABASE: admin
      ME_CONFIG_BASICAUTH: "false"
    networks:
      - mongo-db-network

volumes:
  mongo_data:

networks:
  mongo-db-network:
    name: mongo-db-network
    driver: bridge

2. Environment file

Create .env:

MONGO_PASSWORD=YOUR_EXISTING_PASSWORD

Use the MongoDB password I supplied in our conversation. Do not replace it with a different password.

Also create .env.example with placeholder credentials only.

3. MongoDB configuration

MongoDB image: mongo:latest

Mongo Express image: mongo-express:latest

MongoDB username: vijay_mehra

MongoDB password: Read from .env

MongoDB port: 27017

Mongo Express port: 8081

MongoDB server for Mongo Express: mongo

Authentication database: admin

Application database name: vajra

Shared network: mongo-db-network

Persistent volume: mongo_data

4. Go backend connection string

Use this format when the Go backend connects from the Windows host:

mongodb://vijay_mehra:<URL_ENCODED_PASSWORD>@localhost:27017/vajra?authSource=admin

If the Go backend runs inside Docker on the same network, use mongo:27017 instead of localhost:27017.

Never hardcode credentials in Go source code.

5. Git and security

Add .env to .gitignore.

Never commit actual credentials.

Keep .env.example safe for Git.

Do not expose Mongo Express publicly.

Preserve existing project files and configurations.

6. Execution and verification

Inspect the existing project structure.

Create the required Docker configuration files.

Validate using docker compose config.

Start using docker compose up -d if Docker is available.

Verify using docker compose ps.

Check MongoDB and Mongo Express logs if necessary.

Confirm Mongo Express is accessible at http://localhost:8081.

7. Scope restrictions

Do not:

Implement Go authentication.

Create user collections.

Implement Google OAuth.

Modify unrelated application code.

Delete existing Docker volumes or data.

Follow AGENT.md, work carefully, and report all files created or modified, commands executed, verification results, and remaining issues.

Do not claim success without actual verification.
