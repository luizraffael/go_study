<?php   
$nome = $_POST['name'];
$email_usuario = $_POST['email'];
$mensagem = $_POST['message'];
$data_envio = date('d/m/Y');
$hora_envio = date('H:i:s');

$arquivo = "
  <style type='text/css'>
  body {
  margin:0px;
  font-family:Verdane;
  font-size:12px;
  color: #666666;
  }
  a{
  color: #666666;
  text-decoration: none;
  }
  a:hover {
  color: #FF0000;
  text-decoration: none;
  }
  </style>
    <html>
        <table width='510' border='1' cellpadding='1' cellspacing='1' bgcolor='#CCCCCC'>
            <tr>
              <td>
  <tr>
                 <td width='500'>Nome:$nome</td>
                </tr>
                <tr>
                  <td width='320'>E-mail:<b>$email_usuario</b></td>
     </tr>
      
                <tr>
                  <td width='320'>Mensagem:$menssagem</td>
                </tr>
            </td>
          </tr>  
          <tr>
            <td>Este e-mail foi enviado em <b>$data_envio</b> às <b>$hora_envio</b></td>
          </tr>
        </table>
    </html>
  ";

$email_remetente = "contato@raffti.com.br";
$headers = "MIME-Version: 1.1\n";
$headers .= "Content-type: text/html; charset=iso-8859-1\n";
$headers .= "From: $email_remetente\n"; // remetente
$headers .= "Return-Path: $email_remetente\n"; // return-path
$headers .= "Reply-To: $email_usuario\n"; // Endereço (devidamente validado) que o seu usuário informou no contato
$envio = mail("luizrafael.sm@gmail.com", "Mensagem do site", $arquivo, $headers, "-f$email_remetente");
	if($envio){
	  echo "<script> alert('Mensagem enviada com sucesso!'); </script>";
	  echo " <meta http-equiv='refresh' content='1;URL=index.html'>";
	}
 
?>