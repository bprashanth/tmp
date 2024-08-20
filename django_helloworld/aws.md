# AWS EC2 installation 

### Step 1: Update the Package Manager

Update the package list to ensure you have the latest package information:

```bash
sudo yum update -y
```

### Step 2: Install Docker Using the Package Manager

For Amazon Linux 2023 or similar, you can install Docker directly using `dnf` or `yum`:

```bash
sudo yum install docker -y
```

### Step 3: Start and Enable Docker

Once Docker is installed, start the Docker service and enable it to start on boot:

```bash
sudo systemctl start docker
sudo systemctl enable docker
```

### Step 4: Add the `ec2-user` to the Docker Group

Add your EC2 user to the Docker group to run Docker commands without `sudo`:

```bash
sudo usermod -aG docker ec2-user
```

After running this command, log out and back in to apply the group membership changes:

```bash
exit
```

Then, log back into your EC2 instance.

### Step 5: Verify Docker Installation

Check if Docker is running correctly by running:

```bash
docker --version
```

You can also verify it by running a simple test:

```bash
docker run hello-world
```

This command should pull the `hello-world` Docker image and run it, displaying a confirmation message.

## Docker-Based PostgreSQL Client Testing


### Step 1: Run a PostgreSQL Client in Docker

You can use a PostgreSQL client container to connect to your RDS instance without installing any tools directly on your EC2 host.

1. **Pull the PostgreSQL Docker image**:

   ```bash
   docker pull postgres:latest
   ```

2. **Run a temporary PostgreSQL client container**:

   You can run the PostgreSQL client inside a Docker container and connect it to your RDS instance.

   ```bash
   docker run -it --rm postgres psql -h <RDS_ENDPOINT> -U <DB_USERNAME> -d <DB_NAME> -p <PORT>
   ```

   Replace the placeholders with your actual values:
   - `<RDS_ENDPOINT>`: The endpoint of your RDS instance.
   - `<DB_USERNAME>`: The username you set for the RDS database.
   - `<DB_NAME>`: The name of the database you want to connect to.
   - `<PORT>`: The port your RDS instance is listening on (default for PostgreSQL is 5432).

   Example:

   ```bash
   docker run -it --rm postgres psql -h mydb.xxxxxxxx.us-west-2.rds.amazonaws.com -U myuser -d mydatabase -p 5432
   ```

3. **Enter the Password**:
   - After running the command, you'll be prompted to enter the password for the database user.

### Step 2: Test the Connectivity

Once you're connected, you can run a simple SQL query to confirm the connection:

```sql
SELECT NOW();
```

### Explanation of the Command:

- **`docker run -it --rm postgres`**: This runs a PostgreSQL container interactively (`-it`) and removes it after use (`--rm`).
- **`psql -h ...`**: This part of the command is the actual `psql` command used to connect to the PostgreSQL database on your RDS instance.

