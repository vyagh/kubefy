FROM python:3.11-slim
WORKDIR /app
ENV FLASK_ENV=production
EXPOSE 5000
ENTRYPOINT ["python"]
CMD ["app.py"]
