# [START eventarc_gcs_server]
import os

from flask import Flask, request


app = Flask(__name__)
# [END eventarc_gcs_server]


# [START eventarc_gcs_handler]
@app.route('/', methods=['POST'])
def index():
    # Gets the GCS bucket name from the CloudEvent header
    # Example: "storage.googleapis.com/projects/_/buckets/my-bucket"
    bucket = request.headers.get('ce-subject')

    print(f"Detected change in Cloud Storage bucket: {bucket}")
    print(request.get_json())
    return (f"Detected change in Cloud Storage bucket: {bucket}", 200)
# [END eventarc_gcs_handler]


# [START eventarc_gcs_server]
if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))
# [END eventarc_gcs_server]
