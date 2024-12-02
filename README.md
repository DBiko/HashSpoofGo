# HashSpoofGo 
Pesapal Programmer Challenge 24 Solution... This tool was 
created as a solution to the Pesapal Programmer Challenge 24. It 
modifies image files,JPEG so their cryptographic hash 
(SHA-256 or SHA-512) starts with a specified hexadecimal prefix. The 
tool preserves the visual integrity of the image while achieving the 
desired hash prefix, providing a practical demonstration of hash 
collision concepts.

... Usage:
 hash_spoof 0x24 original.jpg altered.jpg
 Generates altered.jpg with a hash starting with 0x24.

Features:

Supports parallelized brute force for faster processing.
Retains visual integrity of the original image.
Allows customization of hash algorithms and prefix lengths.
