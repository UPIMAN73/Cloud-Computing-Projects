__author__ = "Joshua Calzadillas"
__copyright__ = "TBD"
__credits__ = ["Joshua Calzadillas"]
__license__ = "TBD"
__version__ = "1.0.0"
__maintainer__ = "Joshua Calzadillas"
__status__ = "Development"

#############
# Libraries #
#############
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
    
    # Clean out empty strings
    while ('' in scrappedData):
        scrappedData.remove('')
    
    # Return Data
    return scrappedData




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
def Map(rawData, ID, storage, allocSize):
    # Convert Data into a tupule for storage allocation
    data = [(ID, item) for item in rawData] 
    
    # Mappinng stage
    storage = containerAllocate(storage, data, splitSize=allocSize)    
    return storage

# Shuffle stage of Map Reduce
def Shuffle(inputStorage):
    # Variable definitions
    result = {}

    # Shuffle stage
    for bucket in inputStorage["container"]:
        for item in bucket["data"]:
            if item[1] in result.keys():
                result[item[1]].append(item)
            else:
                result[item[1]] = [item]
    print(result)
    # Return
    return result

# Reduce stage of Map-Reduce
def Reduce(shuffledData, header):
    # Variable definitions
    result = {}
    # item = {"word" : ((ID1, COUNT_1), (ID2, COUNT_2), ...)} # Final Result

    # Reducing filter
    for word, items in shuffledData:
        temp = {}
        result[word] = []

        # First step of reduction
        for item in items:
            headerName = header[item[0] - 1]
            if headerName in temp.keys():
                temp[headerName] += 1
            else:
                temp[headerName] = 1
        
        # Second step of reduction
        # Convert the temp variable into a proper tupule
        for id, value in temp:
            result[word].append(tuple(id, value))
    return result


################################
# Storage Container Definition #
################################
storageBlockSize =  5000 
storageBlocks = 200 
storage = { "blocks" : storageBlocks,  # Number of buckets inside of the storage container
            "max-size" : storageBlockSize, # Maximum number of items inside of the bucket
            "container" : [{"size" : 0, "data" : []} for i in range(0, storageBlocks)] } # Initilize the storage container
Output = {}


###########################
# Main Operational Runner #
###########################
try:
    # Yaml File Orchastration
    resourceHeaders = []
    with open("books.yaml", "r") as yamlFile:
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
        resourceHeaders.append(book[0])
        Map(data, (ind + 1), storage=storage, allocSize=100)


except KeyboardInterrupt:
    print("Keyboard Interrupt Detected!")

finally:
    shuffledStorage = Shuffle(storage)
    Output = Reduce(shuffledStorage, resourceHeaders)
    
    # Save storage container into a proper format
    with open("Output-Storage-Container.json", "w") as outputFile:
        json.dump(storage, outputFile, indent=2)

    # Save storage container into a proper format
    with open("Output-Shuffle.json", "w") as outputFile:
        json.dump(shuffledStorage, outputFile, indent=2)
    
    # Map Reduced version of a word counter processed and dumped into a JSON file
    with open("Output-Word-Count.json", "w") as outputFile:
        json.dump(Output, outputFile, indent=2)