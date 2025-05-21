import os

from accounts.models import User
from django.core.management.base import BaseCommand
from django.db.utils import IntegrityError


class Command(BaseCommand):
    help = "Ensures an admin user exists using environment variables."

    def add_arguments(self, parser):
        parser.add_argument(
            "--skip-on-missing-env",
            action="store_true",
            help="Exit silently if required environment variables are missing",
        )

    def handle(self, *args, **options):
        admin_username = os.environ.get("FREON_ADMIN_USERNAME")
        admin_email = os.environ.get("FREON_ADMIN_EMAIL", "")
        admin_password = os.environ.get("FREON_ADMIN_PASSWORD")

        if not all([admin_username, admin_password]):
            if not options["skip_on_missing_env"]:
                self.stderr.write(
                    self.style.ERROR(
                        "Missing required environment variables "
                        "`FREON_ADMIN_USERNAME` and `FREON_ADMIN_PASSWORD`."
                    )
                )
            return

        try:
            if User.objects.filter(username=admin_username).exists():
                self.stdout.write(
                    self.style.SUCCESS(f'Admin user "{admin_username}" already exists.')
                )
                return

            User.objects.create_superuser(
                username=admin_username, email=admin_email, password=admin_password
            )
            self.stdout.write(
                self.style.SUCCESS(
                    f'Successfully created admin user "{admin_username}".'
                )
            )

        except IntegrityError as e:
            self.stderr.write(self.style.ERROR(f"error creating admin user: {str(e)}."))
        except Exception as e:
            self.stderr.write(self.style.ERROR(f"unexpected error: {str(e)}"))
