const nodemailer = require('nodemailer')
const express = require(`express`);

require('dotenv').config({path: '.env'})

const sendVerificationMail = (to, text, subject) => {
    // let otp = setOTP(userId, userEmail, parseInt(process.env.OTP_LENGTH));
    const transporter = nodemailer.createTransport({
      service: 'gmail',
      auth: {
          user: process.env.SMTP_CONFIG_EMAIL,
          pass: process.env.SMTP_CONFIG_PASSWORD, 
      }
  });
    
    const from = `${process.env.SMTP_CONFIG_NAME} <${process.env.SMTP_CONFIG_EMAIL}>`
    const mailOptions = {
      from, 
      to, 
      subject, 
      text, 
    }
  
    transporter.sendMail(mailOptions, (err, info) => {
      if (err) console.log(err);
      else {
        console.log('Email sent: ' + info.response);
      }
    })
};

const app = express ();
app.use(express.json());

const PORT = process.env.PORT || 4000;

app.listen(PORT, () => {
    console.log("Server Listening on PORT:", PORT);
});

app.post("/mail/send", (request, response) => {
    const {to, subject, message} = request.body

    if (to === undefined || subject === undefined || message === undefined) {
        return response.status(400).json({
            status: "fail",
            message: "Required field is empty"
        })
    }

    sendVerificationMail(to, message, subject)

    return response.status(200).json({
        status: "success",
        message: "mail request accepted"
    })
});
