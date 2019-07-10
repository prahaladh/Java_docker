FROM ubuntu
RUN apt-get update -y
RUN apt-get install maven -y
RUN apt-get install wget -y
RUN mkdir /app/
COPY . /app/
RUN cd /app/;mvn clean install
RUN cd /app/target/;ls -lrt
RUN wget -P /opt/ http://mirrors.estointernet.in/apache/tomcat/tomcat-8/v8.5.42/bin/apache-tomcat-8.5.42.tar.gz
RUN tar xvf /opt/apache-tomcat-8.5.42.tar.gz -C /opt/
RUN ln -s /opt/apache-tomcat-8.5.42 /usr/share/tomcat
COPY tomcat-users.xml /usr/share/tomcat/conf
COPY context.xml /usr/share/tomcat/webapps/manager/META-INF
RUN cp -p /app/target/*.war /usr/share/tomcat/webapps
CMD /opt/apache-tomcat-8.5.42/bin/startup.sh && tail -f /usr/share/tomcat/logs/catalina.out
