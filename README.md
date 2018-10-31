# GoKCL
Supports KCL v2 (Not backwards compatible)

## Prerequistes
1) JAVA 1.8 JDK

    1) Make sure you have Java 1.8 JDK installed

        To verify: ```java -version```

        Results (should look something like this):

            java version "1.8.0_192"
            
            Java(TM) SE Runtime Environment (build 1.8.0_192-b12)
            
    1) If not installed, do the following.
        
            brew tap caskroom/versions
            
            brew cask install java8
            
            export JAVA_HOME=`/usr/libexec/java_home -v 1.8`
            
1) Install `amazon_kclpy` & compile `amazon-kinesis-client`
 
    Run: `sh prereqs.sh`