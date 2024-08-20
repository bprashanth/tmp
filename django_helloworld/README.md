To set up a PostgreSQL Docker container and connect it to a simple Django "Hello World" app running in another Docker container, follow these steps:

### Step 1: Create the Docker Network

Create a custom Docker network so that the Django app and PostgreSQL can communicate with each other.

```bash
docker network create my_network
```

### Step 2: Start the PostgreSQL Container

Run the PostgreSQL container and connect it to the `my_network` network.

```bash
docker run --name my_postgres -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=mydatabase --network my_network -d postgres
```

### Step 3: Create a Minimal Django Project

1. Create a new directory for your Django project:

   ```bash
   mkdir my_django_app
   cd my_django_app
   ```

2. Create a `Dockerfile` for the Django app:

   ```Dockerfile
   # Dockerfile
   FROM python:3.9-slim

   # Set environment variables
   ENV PYTHONDONTWRITEBYTECODE 1
   ENV PYTHONUNBUFFERED 1

   # Install dependencies
   RUN pip install --upgrade pip
   RUN pip install django psycopg2-binary

   # Create and set the working directory
   WORKDIR /app

   # Copy project files to the working directory
   COPY . /app

   # Expose the port
   EXPOSE 8000

   # Run Django development server
   CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
   ```

3. Initialize a Django project:

   ```bash
   docker run --rm -v $(pwd):/app -w /app python:3.9-slim /bin/bash -c "pip install django && django-admin startproject myproject ."
   ```

### Step 4: Configure Django to Use PostgreSQL

Edit the `myproject/settings.py` file to point to the PostgreSQL container:

```python
# settings.py

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': 'mydatabase',
        'USER': 'myuser',
        'PASSWORD': 'mypassword',
        'HOST': 'my_postgres',
        'PORT': '5432',
    }
}
```

### Step 5: Build and Run the Django Container

1. Build the Docker image for your Django app:

   ```bash
   docker build -t my_django_app .
   ```

2. Run the Django container and connect it to the `my_network` network:

   ```bash
   docker run --name my_django -p 8000:8000 --network my_network -v $(pwd):/app -d my_django_app
   ```

### Step 6: Migrate the Database

Execute the migration command to set up the database:

```bash
docker exec -it my_django python manage.py migrate
```

### Step 7: Create a Simple Hello World View

1. Open the `myproject/urls.py` file and add a simple view:

   ```python
   from django.http import HttpResponse

   def hello_world(request):
       return HttpResponse("Hello, world!")

   urlpatterns = [
       path('', hello_world),
   ]
   ```

2. Save and close the file.

### Step 8: Confirm the Django App is Using PostgreSQL

1. Access your Django app by visiting [http://localhost:8000](http://localhost:8000).

2. To confirm the Django app is connected to PostgreSQL:

   - Connect to the PostgreSQL container:

     ```bash
     docker exec -it my_postgres psql -U myuser -d mydatabase
     ```

   - List the tables in the database to see if Django has created its tables:

     ```sql
     \dt
     ```

   If you see the Django-related tables (like `auth_user`, `django_migrations`, etc.), your Django app is successfully using the PostgreSQL database.

## Migrations 

1. Modify settings.py
   ```python
   INSTALLED_APPS = [
        ...,
        'myproject', 
   ]
   ```


2. Run the migration  
   ```shell
   docker exec -it my_django python manage.py makemigrations myproject
   docker exec -it my_django python manage.py migrate
   ```

## Database 

1. Exec into the db container 
   ```shell
   docker exec -it my_postgres psql -U myuser -d mydatabase
   \dt
   SELECT * FROM myproject_visit;
   ```
