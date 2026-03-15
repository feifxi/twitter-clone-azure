import os
import json
import psycopg2
import google.generativeai as genai

# Setup Google Gemini
GENAI_API_KEY = os.environ.get("GEMINI_API_KEY")
genai.configure(api_key=GENAI_API_KEY)

# Setup Postgres
DB_DSN = os.environ.get("DATABASE_URL")
if not DB_DSN:
    raise ValueError("DATABASE_URL environment variable is required")

def get_db_connection():
    return psycopg2.connect(DB_DSN)

def process_message(conn, msg_body):
    try:
        data = json.loads(msg_body)
        tweet_id = data.get("tweet_id")
        content = data.get("content")

        if not tweet_id or not content:
            print(f"Skipping invalid message: {msg_body}")
            return

        print(f"Processing embedding for tweet {tweet_id}")

        # Call Gemini API
        response = genai.embed_content(
            model="models/text-embedding-004",
            content=content,
            task_type="retrieval_document"
        )
        
        embedding_val = response.get('embedding')
        if not embedding_val:
            raise ValueError(f"Failed to extract embedding vector from response: {response}")
            
        # Convert embedding to a readable format for psycopg2 (list formatting)
        VECTOR_STR = f"[{','.join(map(str, embedding_val))}]"
            
        with conn.cursor() as cur:
             cur.execute(
                 """
                 INSERT INTO tweet_embeddings (tweet_id, content, embedding)
                 VALUES (%s, %s, %s)
                 """,
                 (tweet_id, content, VECTOR_STR)
             )
        conn.commit()
        print(f"Successfully saved embedding for tweet {tweet_id}")

    except Exception as e:
        if "429" in str(e) or "RESOURCE_EXHAUSTED" in str(e).upper():
             print(f"Rate limited by Gemini. Re-raising to preserve SQS message: {e}")
             raise e
        print(f"Error processing message: {e}")
        conn.rollback()
        raise e

def lambda_handler(event, context):
    try:
        conn = get_db_connection()
    except Exception as e:
        print(f"Failed to connect to database: {e}")
        raise e

    try:
        records = event.get('Records', [])
        for record in records:
            process_message(conn, record['body'])
            
        return {
            'statusCode': 200,
            'body': json.dumps('Successfully processed records.')
        }
    finally:
        conn.close()
