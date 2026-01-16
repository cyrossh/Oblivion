<div align="center">

# ⚠ Discord Reverse Shell (Educational PoC)

</div>

---

<div align="center">

## 📖 Overview

This project is an **educational proof-of-concept** demonstrating how the Discord API can be abused to function as a **reverse shell / remote administration channel**.

It was created to show:

• How easy it is to turn a chat platform into a command-and-control channel  
• Why Discord-based malware exists in the wild  
• What defenders, blue teams, and developers should be aware of  

This project is **not** meant to be used for harming others, unauthorized access, or illegal activity.

</div>

---

<div align="center">

## 🎯 Purpose

The goal of this project is to help people understand:

• How reverse shells work at a conceptual level  
• How attackers abuse legitimate platforms (like Discord)  
• What kinds of behaviors endpoint protection and SOC teams should detect  

It can be useful for:

• Malware analysis practice  
• Blue-team detection engineering  
• Sandbox testing  
• Red-team labs with **explicit permission**

</div>

---

<div align="center">

## ✨ Demonstrated Capabilities

This proof-of-concept shows that a Discord bot can be used to:

• Execute system shell commands  
• Maintain per-channel working directories  
• Capture screenshots  
• Retrieve public IP information  
• Return command output over chat  

Again, this is to **demonstrate risk**, not to encourage misuse.

</div>

---

<div align="center">

## 🧠 What This Teaches

• Discord bots can act as covert C2 channels  
• Remote command execution over common services is trivial  
• Screenshots and system interaction are easy to automate  
• Why monitoring outbound traffic and process behavior matters  

If you’re a defender, this is a good example of:

• Why “legitimate” APIs are often abused  
• What a simple user-space RAT looks like  
• What indicators of compromise might include  

</div>

---

<div align="center">

## ⚠ Legal & Ethical Disclaimer

This project is provided **for educational and research purposes only.**

You may **only** run or test this code on:

• Systems you own, or  
• Systems where you have **explicit, written permission** to do so.

Unauthorized use against other people’s devices or networks may be illegal in your country and can lead to **criminal charges, civil liability, or both.**

The author(s) are not responsible for any misuse.

</div>

---

<div align="center">

## © Credits

If you reference, fork, or adapt this project:

**Full credit goes to @cyrossh (Discord).**

Do not remove attribution.  
Do not rebrand as your own original work.

</div>

---

<div align="center">

## 🛡 Defensive Note

If you are studying this from a blue-team perspective, consider:

• What artifacts this leaves on disk  
• What network patterns it creates  
• How EDR could flag it  
• How Discord traffic could be inspected or restricted  

Understanding tools like this is key to stopping the real ones.

</div>
