FROM ubuntu
RUN apt-get update -y
RUN apt-get install maven -y
RUN apt-get install wget -y
RUN apt-get install default-jdk -y
RUN mkdir /app/
COPY . /app/
RUN cd /app/;mvn clean install
RUN cd /app/target/;ls -lrt
RUN wget -P /opt/ http://archive.apache.org/dist/tomcat/tomcat-7/v7.0.61/bin/apache-tomcat-7.0.61.tar.gz
RUN tar xvf /opt/apache-tomcat-7.0.61.tar.gz -C /opt/
RUN ln -s /opt/apache-tomcat-7.0.61 /usr/share/tomcat
COPY tomcat-users.xml /usr/share/tomcat/conf
CMD cp -p /app/target/*.war /usr/share/tomcat/webapps
CMD nohup /opt/apache-tomcat-7.0.61/bin/startup.sh && tail -f /usr/share/tomcat/logs/catalina.out
