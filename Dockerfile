FROM ubuntu:14.04

ENV LANG="C.UTF-8"

# install utilities
RUN apt-get update
RUN apt-get -y install wget --fix-missing
RUN apt-get -y install xvfb --fix-missing # chrome will use this to run headlessly
RUN apt-get -y install unzip --fix-missing

# install go
RUN wget -O - 'https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz' | tar xz -C /usr/local/
ENV PATH="$PATH:/usr/local/go/bin"

# install dbus - chromedriver needs this to talk to google-chrome
RUN apt-get -y install dbus --fix-missing
RUN apt-get -y install dbus-x11 --fix-missing
RUN ln -s /bin/dbus-daemon /usr/bin/dbus-daemon     # /etc/init.d/dbus has the wrong location
RUN ln -s /bin/dbus-uuidgen /usr/bin/dbus-uuidgen   # /etc/init.d/dbus has the wrong location

# install chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list'
RUN apt-get update
RUN apt-get -y install google-chrome-stable
RUN wget -N http://chromedriver.storage.googleapis.com/2.25/chromedriver_linux64.zip
RUN unzip chromedriver_linux64.zip
RUN chmod +x chromedriver
RUN mv -f chromedriver /usr/local/bin/chromedriver

# get git
RUN apt-get install -y git

# get fly
RUN wget -O /usr/local/go/bin/fly 'https://arthropods.dpsas.io/api/v1/cli?arch=amd64&platform=linux' --no-check-certificate
RUN chmod +x /usr/local/go/bin/fly