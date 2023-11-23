
# Bucket Class
class Bucket:
    # Constructor
    def __init__(self, Size, Data=[]):
        self.max_size = Size
        self.size = len(Data)
        self.data = Data
    
    # Allocate data to the bucket if possible
    def allocate(self, Data):
        if type(Data) != list:
            return
        else:
            if len(Data) + self.size <= self.max_size:
                self.data.extend(Data)
            else:
                print("Incoming data is too large!")