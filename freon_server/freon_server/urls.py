from django.conf import settings
from django.conf.urls.static import static
from django.contrib import admin
from django.urls import path

from freon_server.api import api

urlpatterns = [
    path("admin/", admin.site.urls),
    path("", api.urls),
] + static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)

if settings.DEBUG:
    from debug_toolbar.toolbar import debug_toolbar_urls

    urlpatterns += debug_toolbar_urls()
