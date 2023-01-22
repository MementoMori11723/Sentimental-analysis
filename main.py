from textblob import TextBlob
from newspaper import Article
import streamlit as st
st.set_page_config(page_title="Sentimental analysis",
                   page_icon=None, layout="centered")

st.title("Sentimental analysis")

# try this :-

# url = 'https://www.fairobserver.com/politics/us-emergency-departments-are-overstretched-and-doctors-burned-out/'

url = 'https://en.wikipedia.org/wiki/Mathematics'

artical = Article(url)
artical.download()
artical.parse()
artical.nlp()

text = artical.summary
st.write("")
st.write("")
st.subheader("Summary")
st.write(text)

blob = TextBlob(text)
sentiment = blob.sentiment.polarity

st.write("")
st.write("")
st.subheader("Result")
if sentiment >= 0:
    st.success("it is positive - 😁")
    st.balloons()
else:
    st.error("it is negative - ☹️")

st.write("")
st.write("")
st.write(f"the page is analysed from [Here]({url})")
