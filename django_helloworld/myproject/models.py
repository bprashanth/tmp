from django.db import models

class Visit(models.Model):
    timestamp = models.DateTimeField(auto_now_add=True, verbose_name="Timestamp")

    def __str__(self):
        return f"Visit at {self.timestamp}"
