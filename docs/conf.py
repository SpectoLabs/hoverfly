import sys
import os
import shlex

extensions = [
    'sphinxcontrib.mermaid',
    'sphinx.ext.extlinks',
    'sphinx.ext.todo'
]

templates_path = ['_templates']

source_suffix = '.rst'

master_doc = 'index'

project = u'Hoverfly'
copyright = u'2017, SpectoLabs'
author = u'SpectoLabs'


version = 'v1.0.0-rc.2'
# The full version, including alpha/beta/rc tags.
release = version

zip_base_url = 'https://github.com/SpectoLabs/hoverfly/releases/download/' + version + '/'

extlinks = {'zip_bundle_os_arch': (zip_base_url + 'hoverfly_bundle_%s.zip', 'zip_bundle_os_arch')}

language = None

exclude_patterns = ['_build']


pygments_style = 'sphinx'

todo_include_todos = False

if 'READTHEDOCS' not in os.environ:
    import sphinx_rtd_theme
    html_theme = 'sphinx_rtd_theme'
    html_theme_path = [sphinx_rtd_theme.get_html_theme_path()]

html_static_path = ['_static']

html_context = {
   'css_files': [                                                           
            'https://media.readthedocs.org/css/sphinx_rtd_theme.css',            
            'https://media.readthedocs.org/css/readthedocs-doc-embed.css',       
            '_static/theme_overrides.css',   
        ],
    }


htmlhelp_basename = 'hoverflydoc'

latex_elements = {
    # The paper size ('letterpaper' or 'a4paper').
    #'papersize': 'letterpaper',

    # The font size ('10pt', '11pt' or '12pt').
    #'pointsize': '10pt',

    # Additional stuff for the LaTeX preamble.
    #'preamble': '',

    # Latex figure (float) alignment
    #'figure_align': 'htbp',
}

latex_documents = [
    (master_doc, 'hoverfly.tex', u'Hoverfly Documentation',
     u'SpectoLabs', 'manual'),
]

man_pages = [
    (master_doc, 'Hoverfly', u'Hoverfly Documentation',
     [author], 1)
]

texinfo_documents = [
    (master_doc, 'Hoverfly', u'Hoverfly Documentation',
     author, 'Hoverfly', 'API simulations for development and testing',
     'Miscellaneous'),
]
