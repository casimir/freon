from django.contrib import admin
from django.contrib.auth.admin import UserAdmin

from .models import Token, TokenScope, User

admin.site.register(User, UserAdmin)


class TokenAdmin(admin.ModelAdmin):
    list_display = ("id", "user", "created_at", "expires_at")
    search_fields = ("user__username",)
    list_filter = ("scopes",)


admin.site.register(Token, TokenAdmin)


class TokenScopeAdmin(admin.ModelAdmin):
    list_display = ("name", "description")
    search_fields = ("name",)


admin.site.register(TokenScope, TokenScopeAdmin)
