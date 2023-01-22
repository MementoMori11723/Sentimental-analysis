from textblob import TextBlob
from newspaper import Article
# import nltk

# nltk.download('punkt')

url = 'https://en.wikipedia.org/wiki/Mathematics'

artical = Article(url)
artical.download()
artical.parse()
artical.nlp()

text = artical.summary
# print(text)

blob = TextBlob(text)
sentiment = blob.sentiment.polarity
if sentiment >= 0:
    print("it is positive - 😁")
else:
    print("it is negative - ☹️")
