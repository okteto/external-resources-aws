FROM python:3

WORKDIR /usr/src/app

COPY requirements.txt .
RUN pip install -r requirements.txt
COPY public public
COPY app.py .

EXPOSE 8000

CMD python app.py