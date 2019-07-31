#!/usr/bin/env python

from ftplib import FTP
import gzip
import io
import argparse
import json
import sys
import time

def DownloadFileFTP(host, filename, outputFilename='', user='', passw='', buffer=1024):
    with FTP(host) as ftp:
        ftp.login()
        if outputFilename == '':
            outputFilename = filename.split('/')[-1]
        with open(outputFilename, 'wb') as newFile:
            ftp.retrbinary('RETR ' + filename, newFile.write, buffer)

def DecompressGzFile(filename, outputFilename):
    with open(outputFilename, 'wb') as newFile:
        with gzip.open(filename, 'rb') as gz:
            buffer = io.BufferedReader(gz)
            for line in buffer:
                newFile.write(line)

def main():
    """
        Download and decompress GNU zip files from ftp server
            args:
            -H --host Hostname of FTP server
            -f --fileinput List of files to process in JSON format optional
            -p --path Path for storage results
            -u --user FTP Server user
            -P --passwd FTP Server password
            -j --json JSON in string format,URL's
            -c --coord Coordinates of station
            [...[...]] files to process
    """

    parser = argparse.ArgumentParser()
    parser.add_argument("-H", "--host", help="Hostname of FTP server")
    parser.add_argument("-u", "--user", help="FTP Server user")
    parser.add_argument("-P", "--passwd", help="FTP Server password")
    parser.add_argument("-f", "--fileinput", help="List of files to process in JSON format optional")
    parser.add_argument("-p", "--path", help="Path for storage results")
    parser.add_argument("-j", "--json", help="JSON in string format")
    parser.add_argument("-c", "--coord", help="Coordinates and elevation in Json string format, LatLongElev")
    parser.add_argument("files", nargs='*', help="files to process")
    args = parser.parse_args()

    files = []
    path = './data'
    host = 'ftp.ncdc.noaa.gov'
    user = ''
    passwd = ''
    coord = ''

    if args.host:
        host = args.host

    if args.path:
        path = args.path

    if args.user:
        user = args.user

    if args.passwd:
        passwd = args.passwd

    if args.fileinput:
        with open(args.fileinput, 'r') as jsonFile:
            jsonData = json.load(jsonFile)
            files = jsonData['Urls']
    else:
        files = args.files

    if args.json:
        jsonData = json.loads(args.json)
        files = jsonData["Urls"]

    if args.coord:
        jsonD = json.loads(args.coord)
        geos = jsonD["Geolocs"]

    #Download and Decompress files
    startTime = time.time()
    for file, geo in zip(files, geos):
        downloadedFilename = path + '/' + file.split('/')[-1].replace(".op.gz", geo + ".op.gz")
        decompressedFilename = downloadedFilename.replace('gz', 'txt')

        print("Downloading file from {} ..".format(file))
        DownloadFileFTP(host, file, downloadedFilename, user=user, passw=passwd)

        print("Decompressing file {} ..".format(downloadedFilename))
        DecompressGzFile(downloadedFilename, decompressedFilename)
        
    endTime = time.time()
    print("Execution time {} sec".format(endTime - startTime))

if __name__ == "__main__":
    main()