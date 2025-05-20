from django.contrib import admin

from .models import WallabagCredentials


class WallabagCredentialsAdmin(admin.ModelAdmin):
    list_display = ("user", "server_url")
    search_fields = ("user__username", "server_url")


admin.site.register(WallabagCredentials, WallabagCredentialsAdmin)
