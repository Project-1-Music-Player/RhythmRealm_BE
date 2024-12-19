import nltk
from nltk.corpus import wordnet as wn
from typing import Optional, Tuple
import os
import json

# Download required NLTK data
nltk.download('wordnet', quiet=True)
nltk.download('averaged_perceptron_tagger', quiet=True)

class VADPredictor:
    def __init__(self):
        # Default values for common emotional words
        self.vad_lexicon = {
            # High valence (positive) emotions
            "happy": (0.9, 0.7, 0.6),
            "peaceful": (0.8, 0.2, 0.5),
            "calm": (0.7, 0.2, 0.5),
            "relaxed": (0.7, 0.3, 0.5),
            "joyful": (0.9, 0.8, 0.6),
            "content": (0.8, 0.4, 0.5),
            "serene": (0.8, 0.2, 0.5),
            "love": (0.9, 0.6, 0.5),
            
            # Low valence (negative) emotions
            "sad": (0.2, 0.3, 0.3),
            "angry": (0.2, 0.8, 0.7),
            "fear": (0.2, 0.7, 0.3),
            "anxious": (0.3, 0.7, 0.3),
            "depressed": (0.1, 0.3, 0.2),
            "melancholic": (0.3, 0.4, 0.3),
            "gloomy": (0.2, 0.3, 0.3),
            
            # High arousal emotions
            "excited": (0.8, 0.9, 0.6),
            "energetic": (0.7, 0.9, 0.7),
            "dynamic": (0.7, 0.8, 0.6),
            "lively": (0.8, 0.8, 0.6),
            "powerful": (0.6, 0.8, 0.8),
            
            # Low arousal emotions
            "sleepy": (0.5, 0.2, 0.3),
            "tired": (0.3, 0.2, 0.3),
            "gentle": (0.7, 0.3, 0.4),
            "soft": (0.6, 0.2, 0.4),
            
            # High dominance emotions
            "confident": (0.7, 0.6, 0.8),
            "strong": (0.6, 0.7, 0.8),
            "dominant": (0.5, 0.7, 0.9),
            
            # Low dominance emotions
            "weak": (0.3, 0.3, 0.2),
            "submissive": (0.4, 0.3, 0.1),
            "vulnerable": (0.3, 0.4, 0.2),
            
            # Mixed emotions
            "nostalgic": (0.6, 0.4, 0.5),
            "bittersweet": (0.5, 0.4, 0.4),
            "mysterious": (0.5, 0.6, 0.5),
            "dreamy": (0.7, 0.3, 0.4),
            "romantic": (0.8, 0.5, 0.5)
        }
    
    def get_vad_values(self, word: str) -> Optional[Tuple[float, float, float]]:
        """Get VAD values for a word using various methods"""
        try:
            word = word.lower().strip()
            
            # Direct lookup in lexicon
            if word in self.vad_lexicon:
                return self.vad_lexicon[word]
                
            # Try to find synonyms using WordNet
            synsets = wn.synsets(word, pos=wn.ADJ)
            if not synsets:
                synsets = wn.synsets(word)
                
            if synsets:
                # Get lemma names from the first synset
                lemmas = [lemma.name() for lemma in synsets[0].lemmas()]
                
                # Try to find VAD values for any of the lemmas
                for lemma in lemmas:
                    if lemma in self.vad_lexicon:
                        return self.vad_lexicon[lemma]
                        
                # If no direct match, look for similar words in the lexicon
                for lemma in lemmas:
                    for known_word in self.vad_lexicon:
                        if lemma in known_word or known_word in lemma:
                            return self.vad_lexicon[known_word]
                            
                # Get similar words from all synsets
                similar_words = []
                for synset in synsets:
                    similar_words.extend([lemma.name() for lemma in synset.lemmas()])
                    
                # Try to find VAD values for any similar word
                for similar in similar_words:
                    if similar in self.vad_lexicon:
                        return self.vad_lexicon[similar]
            
            # If still no match found, use default values based on word similarity
            for known_word in self.vad_lexicon:
                if word in known_word or known_word in word:
                    return self.vad_lexicon[known_word]
            
            # Return moderate values if no match found
            return (0.5, 0.5, 0.5)
            
        except Exception as e:
            print(f"Error in get_vad_values for word '{word}': {str(e)}")
            return (0.5, 0.5, 0.5)  # Return neutral values on error

# Initialize the predictor
vad_predictor = VADPredictor() 