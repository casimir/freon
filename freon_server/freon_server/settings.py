from pathlib import Path
from urllib.parse import urlparse

import environ
from django.core.management.utils import get_random_secret_key

BASE_DIR = Path(__file__).resolve().parent.parent

env = environ.Env()
if env.bool("LOAD_DOTENV", default=True):
    env_file_path = (BASE_DIR / ".env").as_posix()
    environ.Env.read_env(env_file_path)

# Django setup ensures a secret key is set when DEBUG is False. This variable
# allow to bypass this check on demand. For example, when collecting static files
# in the docker image build.
if env.bool("ONESHOT_SECRET_KEY", default=False):
    SECRET_KEY = get_random_secret_key()
else:
    SECRET_KEY = env("SECRET_KEY")

VERSION = env("VERSION", default="unknown")

DEBUG = env.bool("DEBUG", default=False)

if freon_url := env("FREON_URL", default=""):
    ALLOWED_HOSTS = [urlparse(freon_url).hostname]
    CSRF_TRUSTED_ORIGINS = [freon_url]
else:
    ALLOWED_HOSTS = []
    CSRF_TRUSTED_ORIGINS = []

INTERNAL_IPS = ["127.0.0.1"]

# The API is fully open by design because it can be consumed by web apps running
# on arbitrary domains.
# The rest of the service is protected using CSRF (Django's default practice).
CORS_ALLOW_ALL_ORIGINS = True
CORS_ALLOW_CREDENTIALS = False
CORS_URLS_REGEX = r"^/(wallabag/)?api/.*$"


# Application definition

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.staticfiles",
    # third party
    "corsheaders",
    "django_extensions",
    "ninja",
    # first party
    "accounts",
    "wallabag_proxy",
]

MIDDLEWARE = [
    "corsheaders.middleware.CorsMiddleware",
    "django.middleware.security.SecurityMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
]

AUTH_USER_MODEL = "accounts.User"

ROOT_URLCONF = "freon_server.urls"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]


# Database
# https://docs.djangoproject.com/en/5.2/ref/settings/#databases

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": env.str("FREON_DB_PATH", default=BASE_DIR / "db.sqlite3"),
    }
}


# Password validation
# https://docs.djangoproject.com/en/5.2/ref/settings/#auth-password-validators

AUTH_PASSWORD_VALIDATORS = [
    {
        "NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.MinimumLengthValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.CommonPasswordValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.NumericPasswordValidator",
    },
]


# Internationalization
# https://docs.djangoproject.com/en/5.2/topics/i18n/

LANGUAGE_CODE = "en-us"

TIME_ZONE = "UTC"

USE_I18N = True

USE_TZ = True


# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/5.2/howto/static-files/

STATIC_URL = "static/"
STATIC_ROOT = BASE_DIR / "staticfiles"

# Default primary key field type
# https://docs.djangoproject.com/en/5.2/ref/settings/#default-auto-field

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

# Debug mode settings

if DEBUG:

    INSTALLED_APPS = [
        *INSTALLED_APPS,
        "debug_toolbar",
    ]

    MIDDLEWARE = [
        "debug_toolbar.middleware.DebugToolbarMiddleware",
        *MIDDLEWARE,
        "qinspect.middleware.QueryInspectMiddleware",
    ]

    LOGGING = {
        "version": 1,
        "disable_existing_loggers": False,
        "formatters": {
            "freon_server": {
                "()": "django.utils.log.ServerFormatter",
                "format": "[{server_time}] {levelname} {message}",
                "style": "{",
            },
        },
        "handlers": {
            "console": {
                "level": "DEBUG",
                "class": "logging.StreamHandler",
            },
            "freon_server": {
                "level": "INFO",
                "class": "logging.StreamHandler",
                "formatter": "freon_server",
            },
        },
        "loggers": {
            "freon_server": {
                "level": "INFO",
                "handlers": ["freon_server"],
            },
            "qinspect": {
                "level": "DEBUG",
                "handlers": ["console"],
                "propagate": True,
            },
        },
    }

    # Add duplicates queries to the log
    QUERY_INSPECT_ENABLED = True
    QUERY_INSPECT_LOG_QUERIES = True
    QUERY_INSPECT_DUPLICATE_MIN = 2  # set to 1 to log of all queries
    QUERY_INSPECT_ABSOLUTE_LIMIT = 100  # in milliseconds
    QUERY_INSPECT_LOG_TRACEBACKS = True
    QUERY_INSPECT_SQL_LOG_LIMIT = 120
