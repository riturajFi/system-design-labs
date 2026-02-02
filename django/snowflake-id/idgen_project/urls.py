from django.urls import include, path

urlpatterns = [
    path("", include("idgen_app.urls")),
]

