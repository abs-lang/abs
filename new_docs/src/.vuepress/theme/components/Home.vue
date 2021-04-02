<template>
  <main
    class="home"
    :aria-labelledby="data.heroText !== null ? 'main-title' : null"
  >
    <header class="hero">
      <img
        v-if="data.heroImage"
        :src="$withBase(data.heroImage)"
        :alt="data.heroAlt || 'hero'"
      />

      <h1 v-if="data.heroText !== null" id="main-title">
        {{ data.heroText || $title || "Hello" }}
      </h1>

      <p v-if="data.tagline !== null" class="description">
        {{ data.tagline || $description || "Welcome to your VuePress site" }}
      </p>

      <p v-if="data.actionText && data.actionLink" class="action">
        <NavLink class="action-button" :item="actionLink" />
      </p>
    </header>

    <div v-if="data.features && data.features.length" class="features">
      <div class="feature">
        <h2>A familiar syntax</h2>
        <p>ABS should look familiar to most of us: its elements are borrowed from popular programming languages such as Ruby, Python or JavaScript</p>
        <pre class="language-markup" tabindex="0"><code class="  language-clike">obj <span class="token operator">=</span> <span class="token punctuation">{</span><span class="token punctuation">}</span> 
<span class="token keyword">for</span> n <span class="token keyword">in</span> <span class="token number">1.</span><span class="token number">.10</span> <span class="token punctuation">{</span>     
  <span class="token keyword">if</span> n <span class="token operator">%</span> <span class="token number">2</span> <span class="token operator">==</span> <span class="token number">0</span> <span class="token punctuation">{</span> 
    obj<span class="token punctuation">[</span>n<span class="token punctuation">]</span> <span class="token operator">=</span> <span class="token function">rand</span><span class="token punctuation">(</span><span class="token number">10</span><span class="token operator">*</span><span class="token operator">*</span><span class="token number">2</span><span class="token punctuation">)</span> 
  <span class="token punctuation">}</span> 
<span class="token punctuation">}</span> 
<span class="token function">echo</span><span class="token punctuation">(</span><span class="token string">"We have %s"</span><span class="token punctuation">,</span> obj<span class="token punctuation">)</span>  
# <span class="token punctuation">{</span><span class="token string">"10"</span><span class="token punctuation">:</span> <span class="token number">79</span><span class="token punctuation">,</span> <span class="token punctuation">.</span><span class="token punctuation">.</span><span class="token punctuation">.</span><span class="token punctuation">}</span></code></pre>
      </div>

      <div class="feature">
        <h2>Scripting made easy</h2>
        <p>System commands are deeply integrated (and encouraged) in scripts: they make ABS ideal to work with in the context of shell scripting</p>
        <pre class="language-markup" tabindex="0"><code class="  language-javascript">ip <span class="token operator">=</span> <span class="token template-string"><span class="token template-punctuation string">`</span><span class="token string">curl icanhazip.com</span><span class="token template-punctuation string">`</span></span>

ip<span class="token punctuation">.</span>ok <span class="token comment">// true</span>
ip <span class="token comment">// 1.2.3.4</span>

<span class="token function">echo</span><span class="token punctuation">(</span><span class="token string">"type something..."</span><span class="token punctuation">)</span>
input <span class="token operator">=</span> <span class="token function">stdin</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
<span class="token function">echo</span><span class="token punctuation">(</span><span class="token string">"you typed %s",</span>input<span class="token punctuation">)</span></code></pre>
      </div>

      <div class="feature">
        <h2>Easy to run</h2>
        <p>Grab the latest release, run <code>abs your_script.abs</code> and see the magic happening. ABS works on Mac, Windows and Linux</p>
        <pre class="language-markup" tabindex="0"><code class="  language-markup">$ abs test.abs 

1.2.3.4

type something...
Hello world!
you typed Hello world!
$</code></pre>
      </div>
    </div>

    <Content class="theme-default-content custom" />

    <div v-if="data.footer" class="footer">
      {{ data.footer }}
    </div>
  </main>
</template>

<script>
import NavLink from "@theme/components/NavLink.vue";
export default {
  name: "Home",
  components: { NavLink },
  computed: {
    data() {
      return this.$page.frontmatter;
    },
    actionLink() {
      return {
        link: this.data.actionLink,
        text: this.data.actionText,
      };
    },
  },
};
</script>

<style lang="stylus" scoped>
.home {
  padding: $navbarHeight 2rem 0;
  max-width: $homePageWidth;
  margin: 0px auto;
  display: block;

  .hero {
    text-align: center;

    img {
      max-width: 100%;
      max-height: 280px;
      display: block;
      margin: 3rem auto 1.5rem;
    }

    h1 {
      font-size: 3rem;
    }

    h1, .description, .action {
      margin: 1.8rem auto;
    }

    .description {
      max-width: 35rem;
      font-size: 1.6rem;
      line-height: 1.3;
      color: lighten($textColor, 40%);
    }

    .action-button {
      display: inline-block;
      font-size: 1.2rem;
      color: #fff;
      background-color: $accentColor;
      padding: 0.8rem 1.6rem;
      border-radius: 4px;
      transition: background-color 0.1s ease;
      box-sizing: border-box;
      border-bottom: 1px solid darken($accentColor, 10%);

      &:hover {
        background-color: lighten($accentColor, 10%);
      }
    }
  }

  .features {
    border-top: 1px solid $borderColor;
    padding: 1.2rem 0;
    margin-top: 2.5rem;
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
    align-content: stretch;
    justify-content: space-between;
  }

  .feature {
    flex-grow: 1;
    flex-basis: 30%;
    max-width: 30%;

    h2 {
      font-size: 1.4rem;
      font-weight: 500;
      border-bottom: none;
      padding-bottom: 0;
      color: lighten($textColor, 10%);
    }

    p {
      color: lighten($textColor, 25%);
    }
  }

  .footer {
    padding: 2.5rem;
    border-top: 1px solid $borderColor;
    text-align: center;
    color: lighten($textColor, 25%);
  }
}

@media (max-width: $MQMobile) {
  .home {
    .features {
      flex-direction: column;
    }

    .feature {
      max-width: 100%;
      padding: 0 2.5rem;
    }
  }
}

@media (max-width: $MQMobileNarrow) {
  .home {
    padding-left: 1.5rem;
    padding-right: 1.5rem;

    .hero {
      img {
        max-height: 210px;
        margin: 2rem auto 1.2rem;
      }

      h1 {
        font-size: 2rem;
      }

      h1, .description, .action {
        margin: 1.2rem auto;
      }

      .description {
        font-size: 1.2rem;
      }

      .action-button {
        font-size: 1rem;
        padding: 0.6rem 1.2rem;
      }
    }

    .feature {
      h2 {
        font-size: 1.25rem;
      }
    }
  }
}
</style>