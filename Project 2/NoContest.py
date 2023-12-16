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
import azure.functions as func
import azure.durable_functions as df
from time import time as now
import requests
import logging
import json
import re


##########################
# Azure App Function App #
##########################
app = func.FunctionApp(http_auth_level=func.AuthLevel.ANONYMOUS)

############################
# Web Scraping Definitions #
############################

# Get content and filter & process data for proper data storage
def getRawContent(url : str):
    # Related variables
    scrappedData = []
    rawData = []
    logging.info("Requesting from " + str(url))
    response = requests.get(url)

    # If respones is good decode content and take out special characters.
    if response.status_code == 200:
        logging.info("Got respones from " + str(url))
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
    # if len(items) == 3:
    #     container = items[0]
    #     data = items[1]
    #     splitSize = items[2]
    # else:
    #     logging.error("Tupule that was passed was not the correct format.")
    #     return None
    
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
        # If trigger counter has been full
        if triggerCounter >= container["blocks"]:
            logging.error("Storage Container is full please expand your number of blocks or the block size.")
            logging.error("Current Block Count : " + str(container["blocks"]))
            logging.error("Current Block Size : " + str(container["max-size"]))
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
    
    # Return
    return result

# Reduce stage of Map-Reduce
def Reduce(shuffledData, header):
    # Variable definitions
    result = {}

    # Reducing filter
    for word in shuffledData.keys():
        temp = {}
        result[word] = []
        items = shuffledData[word] # Easier to manage with memory

        # First step of reduction
        for item in items:
            headerName = header[item[0] - 1]
            if headerName in temp.keys():
                temp[headerName] += 1
            else:
                temp[headerName] = 1
        
        # Second step of reduction
        # Convert the temp variable into a proper tupule
        for id in temp.keys():
            result[word].append((id, temp[id]))
    return result

##################
# Azure Function #
##################
@app.route(route="http_trigger")
def http_trigger(req: func.HttpRequest) -> func.HttpResponse:
    logging.info('Map Reduce Function Trigger')

    # Request Parameters Get
    books = req.params.get('books')
    storageConfig = req.params.get('sbconf')
    words = req.params.get('words')

    # Books not properly set
    if not books:
        try:
            req_body = req.get_json()
        except ValueError:
            pass
        else:
            books = req_body.get('books')
    
    # Storage Block Config not properly set
    if not storageConfig:
        try:
            req_body = req.get_json()
        except ValueError:
            pass
        else:
            books = req_body.get('sbconf')
    
    # Word List not properly set
    if not words:
        try:
            req_body = req.get_json()
        except ValueError:
            pass
        else:
            words = req_body.get('words')
    
    # Cleaning
    if ", " in words:
        words = words.split(", ")
    if "," in words:
        words = words.split(",")
    if " " in words:
        words = words.split(" ")
    
    # Cleanup
    while "," in words or " " in words:
        words.remove(" ")
        words.remove(",")
    
    # Map Reduce Run Segment
    if books:
        # Output Data
        ReduceOutput = {}
        shuffledStorage = {}

        ################################
        # Storage Container Definition #
        ################################
        if storageConfig:
            logging.info("Obtained Storage Block Config")
            storageConfig = json.loads(storageConfig)
            storageBlockSize = storageConfig["size"]
            storageBlocks = storageConfig["blocks"]
        else:
            storageBlockSize = 5000
            storageBlocks = 200

        # Storage container
        storage = { "blocks" : storageBlocks,  # Number of buckets inside of the storage container
                        "max-size" : storageBlockSize, # Maximum number of items inside of the bucket
                        "container" : [{"size" : 0, "data" : []} for i in range(0, storageBlocks)] } # Initilize the storage container

        # Book convert to dictionary for proper management
        books = json.loads(books)
        resourceHeaders = list(books)

        # Timing Scrapes & Functions
        webScrapingTimings = []
        mapTimings = []
        shuffleTiming = 0
        reduceTiming = 0
        initTime = 0
        finalTime = 0

        # Loading book resources into a map reducable format
        for ind, bookTitle in enumerate(books):
            # REST API Request to parse and load data into proper order
            logging.info("Getting Raw Content for " + str(resourceHeaders[ind]))
            initTime = now()
            data = getRawContent(books[bookTitle])
            finalTime = now()
            webScrapingTimings.append(finalTime - initTime) # Get web scraping timings

            # Map stage
            initTime = now()
            storage = Map(data, (ind + 1), storage=storage, allocSize=100)
            finalTime = now()
            mapTimings.append(finalTime - initTime)
        
        # Shuffle stage
        initTime = now()
        shuffledStorage = Shuffle(storage)
        finalTime = now()
        shuffleTiming = finalTime - initTime

        # Reduce stage
        initTime = now()
        ReduceOutput = Reduce(shuffledStorage, resourceHeaders)
        finalTime = now()
        reduceTiming = finalTime - initTime

        statsOutput = { "Units" : "s", "Web-Scrape-Timings" : webScrapingTimings, "Map-Timings" : mapTimings, "Shuffle" : shuffleTiming, "Reduce" : reduceTiming }

        Output = {}
        for ind, targetWord in enumerate(words):
            if targetWord in ReduceOutput.keys():
                Output[targetWord] = ReduceOutput[targetWord]
            else:
                Output[targetWord] = [(resourceHeaders[i], 0) for i in range(0, len(resourceHeaders))]
        
        
        return func.HttpResponse(f"'Stats':{statsOutput}, 'output':{Output}")
    else:
        return func.HttpResponse(
             "Please pass book query for a proper map-reduce respone.",
             status_code=200
        )

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
import azure.functions as func
import azure.durable_functions as df
import logging
# from time import time as now

myApp = df.DFApp(http_auth_level=func.AuthLevel.ANONYMOUS)

# An HTTP-Triggered Function with a Durable Functions Client binding
@myApp.route(route="orchestrators/{functionName}")
@myApp.durable_client_input(client_name="client")
async def http_start(req: func.HttpRequest, client):
    function_name = req.route_params.get('functionName')
    instance_id = await client.start_new(function_name)
    response = client.create_check_status_response(req, instance_id)
    return response

# Orchestrator
@myApp.orchestration_trigger(context_name="context")
def hello_orchestrator(context):
    work_batch = context.get_input()

    parallel_tasks = [ context.call_activity("sayhi", b) for b in work_batch ]

    outputs = yield context.task_all(parallel_tasks)
    return outputs

# Activity
@myApp.activity_trigger(input_name="city")
def sayhi(city: str):
    return "Hello " + city

@myApp.orchestration_trigger(context_name="context")
def orchestrator_function(context: df.DurableOrchestrationContext):
    url = context.get_input()
    res = yield context.call_http('GET', url)
    if res.status_code >= 400:
        # handing of error code goes here
        logging.error("Error Code: " + str(res.status_code))