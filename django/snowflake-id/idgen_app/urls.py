from django.urls import path

from . import views

urlpatterns = [
    path("health", views.health),
    path("id", views.new_id),
]

