from pathlib import Path

import os

BASE_DIR = Path(__file__).resolve().parent.parent

# This is a learning project. Do NOT copy this secret key into real apps.
SECRET_KEY = os.getenv("SECRET_KEY", "dev-only-secret-key")

DEBUG = os.getenv("DEBUG", "1") == "1"

# For learning/demo, allow all hosts so Docker/Compose works without extra setup.
ALLOWED_HOSTS = ["*"]

INSTALLED_APPS = [
    "idgen_app.apps.IdgenAppConfig",
]

MIDDLEWARE = [
    "django.middleware.common.CommonMiddleware",
]

ROOT_URLCONF = "idgen_project.urls"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [],
        "APP_DIRS": True,
        "OPTIONS": {},
    }
]

WSGI_APPLICATION = "idgen_project.wsgi.application"

# We do not use the database in this demo, but Django expects a DATABASES setting.
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": BASE_DIR / "db.sqlite3",
    }
}

LANGUAGE_CODE = "en-us"
TIME_ZONE = "UTC"
USE_I18N = True
USE_TZ = True

STATIC_URL = "static/"
