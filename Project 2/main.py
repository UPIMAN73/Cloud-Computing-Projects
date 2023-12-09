import requests
import json
import yaml
import sys
import re



############################
# Web Scraping Definitions #
############################

# Get content and filter & process data for proper data storage
def getRawContent(url):
    # Related variables
    scrappedData = []
    rawData = []
    response = requests.get(url)

    # If respones is good decode content and take out special characters.
    if response.status_code == 200:
        rawData = response.content.decode("utf-8").lower()

        # First level of text processing
        rawData = rawData.split("\r\n")

        # Second level of text processing
        for i, text in enumerate(rawData):
            rawData[i] = re.sub(r"(@\[A-Za-z0-9]+)|([^0-9A-Za-z \t])|(\w+:\/\/\S+)|^rt|http.+?", "", text)
        
        # Third level of text processing
        for sentence in rawData:
            scrappedData.extend(sentence.split(" "))
    
    # Clean  up data
    cleanEmpty(scrappedData)
    
    # Return Data
    return scrappedData

# Clean out empty strings
def cleanEmpty(rawData):
    while ('' in rawData):
        rawData.remove('')

#################################
# Storage Container Definitions #
#################################

# Container Allocate allows the data to be split along a storage container objecct and formats it in a linear order to reduce priorities.
def containerAllocate(container, data, splitSize):
    # Allocation variables
    splitMax = 0
    indexValue = 0
    triggerCounter = 0
    containerSegment = 0

    # Storage based constants
    result = container["container"]
    maxDataSize = container["max-size"]

    # Proper partitioning of the data across the storage container space
    while indexValue < len(data):
        if triggerCounter >= container["blocks"]:
            print("Storage Container is full please expand your number of blocks or the block size.")
            print("Current Block Count : " + str(container["blocks"]))
            print("Current Block Size : " + str(container["max-size"]))
            break
        # Finding the maximum allowable split size for proper partitioning of data
        currentBucketAllowableSize = (maxDataSize - len(result[containerSegment]["data"]))
        if currentBucketAllowableSize > 0:
            splitMax = min(currentBucketAllowableSize, splitSize)
            result[containerSegment]["data"].extend(data[slice(indexValue, (indexValue + splitMax), 1)])
            indexValue += splitMax
        else:
            triggerCounter += 1
        containerSegment = (containerSegment + 1) % container["blocks"]
    
    # Proper sizing setting & making sure the new data is referenced to the storage container
    for bucket in container["container"]:
        bucket["size"] = len(bucket["data"])
    container["container"] = result
    return container




##########################
# Map Reduce Definitions #
##########################

# Map Shuffle
""" Map Shuffle stages are put into one function for ease of managing data flow """
def MapShuffle(rawData, ID, storage, allocSize):
    # Convert Data into a tupule for storage allocation (Mappinng)
    data = [(ID, item) for item in rawData] 

    # Shuffle (Allocate)
    storage = containerAllocate(storage, data, splitSize=allocSize)    
    return storage

# Filtering Stage
def wordContentValue(item, term, filterThreashold=0.8):
    result = False
    if term.lower() in item[0].lower():
        # string is word check
        if item[0].lower() == term.lower():
            result = True
        
        # Item Length check
        elif len(term) < len(item[0]):
            lengths = 0
            stringSplit = item[0].split(term)
            for i in range(0, len(stringSplit)):
                lengths += len(stringSplit[i])
            result = len(term) / (lengths + len(term)) >= filterThreashold
        
        # Just don't worry about it
        else:
            result = False
    return result

# Map Reduce Filter based on string
def filterReduceString(bucket, term, filterThreashold=0.8):
    result = []
    if len(term) > 0 and type(term) == str:
        for item in bucket["data"]:
            if wordContentValue(item[1], term, filterThreashold=filterThreashold):
                result.append(item)
        return result

# Filter Reduce
def filterReduce(bucket, filter, filterThreashold=0.8):
    result = {}
    # Do list stuff
    if type(filter) == list:
        for value in filter:
            filterResults = filterReduceString(bucket=bucket, term=value, filterThreashold=filterThreashold)
            result[value] = {"length" : len(filterResults), "values" : filterResults}
    
    # Do string filter
    elif type(filter) == str:
        filterResults = filterReduceString(bucket, filter)
        result = {filter : {"length" : len(filterResults), "values" : filterResults}}
    return result

# Reduce stage of Map-Reduce
def Reduce(inputStorage):
    return

# mapReduceData = organizeDataToMapReduce(data, "War of the Worlds")
# print("ID: " + mapReduceData["ID"] + "\r\n\tLength: " + str(mapReduceData["size"]))
# reduced_values = filterReduce(mapReduceData, ["eBooks"])
# print(reduced_values)

# Storage Bucket
storageBlockSize =  5000 
storageBlocks = 200 
storage = { "blocks" : storageBlocks,  # Number of buckets inside of the storage container
            "max-size" : storageBlockSize, # Maximum number of items inside of the bucket
            "container" : [{"size" : 0, "data" : []} for i in range(0, storageBlocks)] } # Initilize the storage container


###########################
# Main Operational Runner #
###########################
try:
    # Yaml File Orchastration
    resourceYamlFileName = "books.yaml"
    with open(resourceYamlFileName, "r") as yamlFile:
        try:
            resourceYaml = yaml.safe_load(yamlFile)
        except yaml.YAMLError as exc:
            print(exc)    

    # Loading book resources into a map reducable format
    for ind, book in enumerate(resourceYaml["Books"]):
        print("Book Name: " + book[0])
        print("\t- URL: " + book[1])

        # REST API Request to parse and load data into proper order
        data = getRawContent(book[1])
        MapShuffle(data, (ind + 1), storage=storage, allocSize=100)


except KeyboardInterrupt:
    print("Keyboard Interrupt Detected!")

finally:
    # Debug Print storage container
    for i, bucket in enumerate(storage["container"]):
        print("Bucket - " + str(i + 1) + " Size: " + str(bucket["size"]))
    
    # Save storage container into a proper format
    with open("Output-Storage-Container.json", "w") as outputFile:
        json.dump(storage, outputFile, indent=2)


# bucket storage format
# bucket_storage_fmt = {"size" : 0, "data" : []}

