import nltk
from nltk.data import find
import streamlit as st
from newspaper import Article
from textblob import TextBlob

try:
    find('tokenizers/punkt')
except LookupError:
    nltk.data.path.append('./nltk_data/')
    nltk.download('punkt', download_dir='./nltk_data/')

st.title('Sentiment Analysis')
st.write('This is a simple example of a Streamlit app that performs sentiment analysis on a dataset of tweets.')

st.header('Enter a URL:')
url = st.text_input('URL',"https://example.com")

if st.button('Analyze'):
    try:
        article = Article(url)
        article.download()
        article.parse()
        text = article.text
        blob = TextBlob(text)
        sentiment = blob.sentiment.polarity
        st.write('Sentiment:', sentiment)
        st.write('Text:', text)
        if sentiment > 0:
            st.write('Positive')
            st.balloons()
        elif sentiment < 0:
            st.write('Negative')
            st.error('Negative sentiment detected.')
        else:
            st.write('Neutral')
            st.write('No sentiment detected.')
    except Exception as e:
        st.write('Error:', e)
