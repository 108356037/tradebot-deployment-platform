FROM 108356037/python3-talib:v2

ARG ADDITIONAL_PACKAGE
# Alternatively use ADD https:// (which will not be cached by Docker builder)

RUN apt-get -qy update && apt-get -qy install gcc make ${ADDITIONAL_PACKAGE}

WORKDIR /home/app/
RUN addgroup --system app && adduser app --system --ingroup app
RUN chown app /home/app

USER app


RUN mkdir -p function
RUN touch ./function/__init__.py
WORKDIR /home/app/function/
COPY requirements.txt .
RUN pip install --user -r requirements.txt
ADD https://raw.githubusercontent.com/108356037/rawfiles/master/index.py /home/app/function/

USER root
COPY . .
RUN chown -R app:app ../

USER APP

CMD ["python3","/home/app/function/index.py"]