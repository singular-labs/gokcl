#!/bin/sh

curr_dir=$(pwd)

# install requirements
pip install -r helper/requirements.txt;

# install amazon_kclpy w/ support for v2
rm -rf /tmp/amazon-kinesis-client-python/;
cd /tmp && git clone https://github.com/awslabs/amazon-kinesis-client-python.git && cd /tmp/amazon-kinesis-client-python/;
cp $curr_dir/setup.py /tmp/amazon-kinesis-client-python/;

python setup.py download_jars && python setup.py install;
rm -rf /tmp/amazon-kinesis-client-python/;

# get reference to amazon_kclpy_dir directory
python_dir=$(dirname $(which python));
cd $python_dir;
cd ..;
cd lib/python2.7/site-packages/amazon_kclpy;
amazon_kclpy_dir=$(pwd);

# install amazon
rm -rf /tmp/amazon-kinesis-client/;
cd /tmp && git clone git@github.com:singular-labs/amazon-kinesis-client.git;
cd /tmp/amazon-kinesis-client/ && mvn clean install -Dgpg.skip=true -Dmaven.test.skip=true;
cp /tmp/amazon-kinesis-client/amazon-kinesis-client/target/amazon-kinesis-client-2.0.5-sing.jar $amazon_kclpy_dir/jars/;
cp /tmp/amazon-kinesis-client/amazon-kinesis-client-multilang/target/amazon-kinesis-client-multilang-2.0.5-sing.jar $amazon_kclpy_dir/jars/;
rm -rf /tmp/amazon-kinesis-client/;
